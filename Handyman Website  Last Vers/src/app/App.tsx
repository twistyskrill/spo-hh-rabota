import React, { useEffect, useState } from 'react';
import { Layout } from './components/Layout';
import { AuthPage } from './components/AuthPage';
import { HomePage } from './components/HomePage';
import { ProfilePage } from './components/ProfilePage';
import { ChatPage } from './components/ChatPage';
import { UserProfilePage } from './components/UserProfilePage';
import { UserProfileEditor } from './components/UserProfileEditor.tsx';
import { HandymanProfileEditor } from './components/HandymanProfileEditor';
import { Handyman, Announcement } from './data';
import { AdminPanel } from './components/AdminPanel';
import { api } from './api';

type Page = 'auth' | 'home' | 'profile' | 'chat' | 'profile-user' | 'admin';
type UserRole = 'user' | 'handyman' | null;

type UserProfile = {
  name: string;
  email: string;
  phone: string;
  location: string;
  memberSince: string;
  worker?: any;
};

export default function App() {
  const [currentPage, setCurrentPage] = useState<Page>('auth');
  const [userRole, setUserRole] = useState<UserRole>(null);
  const [isAdmin, setIsAdmin] = useState(false);
  const [selectedHandymanId, setSelectedHandymanId] = useState<string | null>(null);
  const [currentUserProfile, setCurrentUserProfile] = useState<UserProfile | null>(null);
  const [showCreateUserProfile, setShowCreateUserProfile] = useState(false);
  const [isEditingUserProfile, setIsEditingUserProfile] = useState(false);
  const [currentHandymanId, setCurrentHandymanId] = useState<string | null>(null);
  const [showCreateHandyman, setShowCreateHandyman] = useState(false);
  const [isEditingOwnHandymanProfile, setIsEditingOwnHandymanProfile] = useState(false);
  
  // Restore auth state on refresh if token exists
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) return;

    const restoreSession = async () => {
      try {
        const profile = await api.getProfile();
        const savedRole = localStorage.getItem('preferredRole') as UserRole;
        const role: UserRole = savedRole || (profile?.have_worker_profile || profile?.worker ? 'handyman' : 'user');
        localStorage.setItem('preferredRole', role as string);
        setCurrentUserProfile({
          name: profile.name || '',
          email: profile.email || '',
          phone: profile.phone || '',
          location: profile.worker?.location || '',
          memberSince: new Date().toLocaleDateString(),
          worker: profile.worker,
        });
        if (profile.role === 'admin' || profile.role === 'administrator') {
          setIsAdmin(true);
        }
        setUserRole(role);
        if (role === 'handyman' && profile.worker) {
           setCurrentHandymanId(String(profile.id));
        }
        setCurrentPage('home');
      } catch (error) {
        console.error('Failed to restore session from token', error);
        localStorage.removeItem('token');
      }
    };

    restoreSession();
  }, []);

  // Global State
  const [handymen, setHandymen] = useState<Handyman[]>([]);
  const [announcements, setAnnouncements] = useState<Announcement[]>([]);
  const [pendingAnnouncements, setPendingAnnouncements] = useState<Announcement[]>([]);
  const [pendingHandymen, setPendingHandymen] = useState<Handyman[]>([]);

  // Initial load from API (public ads list + handymen)
  useEffect(() => {
    const loadInitialData = async () => {
      try {
        const [adsRes, handymenRes] = await Promise.all([
          api.getAds(20, 0),        // GET /ads
          api.getHandymen(20, 0),   // GET /handyman
        ]);

        // Объявления: adsRes: { ads: [...], total, limit, offset }
        const mappedAnnouncements: Announcement[] = (adsRes.ads || []).map((ad: any) => ({
          id: ad.id,
          title: ad.title,
          category: ad.category_name || ad.category?.name || 'Без категории',
          status: 'Открыто', // публичная лента возвращает только approved
          date: ad.created_at ? new Date(ad.created_at).toLocaleDateString() : '',
          budget: `${ad.price} руб.`,
          handyman: ad.user_name || ad.user?.name || null,
          location: ad.location,
          description: '', // описание в модели не предусмотрено
        }));
        setAnnouncements(mappedAnnouncements);

        // Мастера: handymenRes: { workers: [...], pagination: {...} }
        const mappedHandymen: Handyman[] = (handymenRes.workers || []).map((w: any) => ({
          id: String(w.id),
          name: w.name,
          skill: (w.categories && w.categories[0]?.name) || 'Мастер',
          rating: 0, // в API пока нет рейтинга
          hourlyRate: w.hourly_rate || 0,
          description: w.description || '',
          reviews: [],
          email: w.email,
        }));
        setHandymen(mappedHandymen);
      } catch (error) {
        console.error('Failed to load data from API', error);
      }
    };
    loadInitialData();
  }, []);

  // Load pending items for admin moderation when admin panel is opened
  useEffect(() => {
    const loadPendingForAdmin = async () => {
      if (!isAdmin || currentPage !== 'admin') return;
      try {
        const [adsRes, workersRes] = await Promise.all([
          api.getAdminAds('pending'),
          api.getAdminWorkers('pending'),
        ]);

        const mappedPendingAnnouncements: Announcement[] = (adsRes.ads || []).map((ad: any) => ({
          id: ad.id,
          title: ad.title,
          category: ad.category_name || 'Без категории',
          status: ad.status || 'pending',
          date: ad.created_at ? new Date(ad.created_at).toLocaleDateString() : '',
          budget: `${ad.price} руб.`,
          handyman: null,
          location: ad.location,
          description: '', // в админском списке достаточно заголовка
        }));
        setPendingAnnouncements(mappedPendingAnnouncements);

        const mappedPendingHandymen: Handyman[] = (workersRes.workers || []).map((w: any) => ({
          id: String(w.user_id || w.id),
          name: w.name,
          skill: (w.categories && w.categories[0]?.name) || 'Мастер',
          rating: 0,
          hourlyRate: 0,
          description: w.description || '',
          reviews: [],
          email: w.email,
        }));
        setPendingHandymen(mappedPendingHandymen);
      } catch (error) {
        console.error('Failed to load pending items for admin', error);
      }
    };

    loadPendingForAdmin();
  }, [isAdmin, currentPage]);

  const handleLogin = async (role: 'user' | 'handyman', userData?: any) => {
    try {
      const profile = await api.getProfile();
      localStorage.setItem('preferredRole', role);
      setUserRole(role);

      setCurrentUserProfile({
        name: profile.name || '',
        email: profile.email || '',
        phone: profile.phone || '',
        location: profile.worker?.location || '',
        memberSince: new Date().toLocaleDateString(),
        worker: profile.worker,
      });

      if (profile.role === 'admin' || profile.role === 'administrator') {
        setIsAdmin(true);
      }
      
      if (role === 'handyman' && profile.worker) {
         setCurrentHandymanId(String(profile.id));
      } else {
         setCurrentHandymanId(null);
      }
      setCurrentPage('home');
    } catch (e) {
      console.error('Failed to fetch profile right after login', e);
      // Fallback
      setUserRole(role);
      setCurrentPage('home');
    }
  };

  const handleNavigate = (page: string) => {
      if (page === 'home') {
          setCurrentPage('home');
          setSelectedHandymanId(null);
      } else if (page === 'profile-user') {
          // If role is handyman, maybe they have a different profile page?
          // For now, let's assume both go to the same UserProfile structure 
          // or we restrict it. The prompt didn't specify a separate profile page for Handymen 
          // other than the public one.
          // Let's keep it simple: Handymen can also post jobs if they want, or we can block it.
          // But likely a Handyman Dashboard would be better.
          // For this specific request, we'll just navigate to the standard profile page.
          setCurrentPage('profile-user');
          setSelectedHandymanId(null);
      }
  };

  const handleSelectHandyman = async (id: string) => {
    setSelectedHandymanId(id);
    setCurrentPage('profile');
    try {
      const data = await api.getHandymanById(Number(id));
      if (data) {
        // Find existing to preserve reviews or local state if desired, or just overwrite with API data
        const apiHandyman: Handyman = {
          id: String(data.user_id || data.id || id),
          name: data.name || 'Мастер',
          skill: data.categories?.[0]?.name || 'Специалист',
          rating: 0,
          hourlyRate: 0,
          description: data.description || '',
          reviews: [],
          email: data.email || ''
        };
        // Update handymen list with the detailed data
        setHandymen(prev => {
          const exists = prev.find(h => h.id === apiHandyman.id);
          if (exists) {
            return prev.map(h => h.id === apiHandyman.id ? { ...h, ...apiHandyman, reviews: h.reviews } : h);
          }
          return [...prev, apiHandyman];
        });
      }
    } catch (e) {
      console.error('Failed to load handyman details', e);
    }
  };

  const handleChat = (id: string) => {
    setSelectedHandymanId(id);
    setCurrentPage('chat');
  };

  const handleBackToHome = () => {
    setCurrentPage('home');
    setSelectedHandymanId(null);
  };

  const handleBackToProfile = () => {
    setCurrentPage('profile');
  };

  const handleLogout = () => {
    setCurrentPage('auth');
    setUserRole(null);
    setIsAdmin(false);
    setSelectedHandymanId(null);
    setCurrentUserProfile(null);
    setShowCreateUserProfile(false);
    setIsEditingUserProfile(false);
    setCurrentHandymanId(null);
    setShowCreateHandyman(false);
    setIsEditingOwnHandymanProfile(false);
    localStorage.removeItem('token');
  };

  const currentOwnHandymanProfile = userRole === 'handyman' && currentUserProfile
    ? {
        id: currentHandymanId || Date.now().toString(),
        name: currentUserProfile.name,
        skill: (currentUserProfile as any).worker?.specialization?.[0]?.name || 'Мастер',
        rating: 0,
        hourlyRate: 0,
        description: (currentUserProfile as any).worker?.description || '',
        reviews: [],
        email: currentUserProfile.email,
        location: (currentUserProfile as any).worker?.location || '',
      } as Handyman
    : null;

  // State Updaters
  const handleAddReview = (handymanId: string, review: any) => {
    setHandymen(prev => prev.map(h => {
      if (h.id === handymanId) {
        const newReviews = [review, ...h.reviews];
        const totalRating = newReviews.reduce((acc, r) => acc + r.rating, 0);
        const newRating = Number((totalRating / newReviews.length).toFixed(1));
        return { ...h, reviews: newReviews, rating: newRating };
      }
      return h;
    }));
  };

  const handleAddAnnouncement = (announcement: Announcement) => {
    setPendingAnnouncements(prev => [announcement, ...prev]);
  };

  const handleApproveAnnouncement = (announcementId: number) => {
    setPendingAnnouncements(prev => {
      const announcementToApprove = prev.find(a => a.id === announcementId);
      if (announcementToApprove) {
        // шлем запрос на бэкенд, но не блокируем UI
        api.approveAd(announcementId).catch(err =>
          console.error('Failed to approve ad', err),
        );
        const approvedAnnouncement = {
          ...announcementToApprove,
          status: 'Открыто',
        };
        setAnnouncements(current => [approvedAnnouncement, ...current]);
      }
      return prev.filter(a => a.id !== announcementId);
    });
  };

  const handleRejectAnnouncement = (announcementId: number) => {
    setPendingAnnouncements(prev => {
      const announcementToReject = prev.find(a => a.id === announcementId);
      if (announcementToReject) {
        api.rejectAd(announcementId).catch(err =>
          console.error('Failed to reject ad', err),
        );
      }
      return prev.filter(a => a.id !== announcementId);
    });
  };

  const handleApproveHandyman = (handymanId: string) => {
    setPendingHandymen(prev => {
      const handymanToApprove = prev.find(h => h.id === handymanId);
      if (handymanToApprove) {
        api.approveWorker(Number(handymanId)).catch(err =>
          console.error('Failed to approve worker', err),
        );
        setHandymen(current => [handymanToApprove, ...current]);
      }
      return prev.filter(h => h.id !== handymanId);
    });
  };

  const handleRejectHandyman = (handymanId: string) => {
    setPendingHandymen(prev => {
      const handymanToReject = prev.find(h => h.id === handymanId);
      if (handymanToReject) {
        api.rejectWorker(Number(handymanId)).catch(err =>
          console.error('Failed to reject worker', err),
        );
      }
      return prev.filter(h => h.id !== handymanId);
    });
  };

  if (currentPage === 'auth') {
    return <AuthPage onLogin={(role, userData) => {
      if (userData?.role === 'admin' || userData?.role === 'administrator') {
        setIsAdmin(true);
      }
      handleLogin(role, userData);
    }} />;
  }

  return (
    <Layout
      onNavigate={handleNavigate}
      currentPage={currentPage}
      onLogout={handleLogout}
      isAdmin={isAdmin}
      onOpenAdminPanel={() => setCurrentPage('admin')}
    >
      {currentPage === 'admin' && (
        <AdminPanel
          pendingAnnouncements={pendingAnnouncements}
          pendingHandymen={pendingHandymen}
          onApproveAnnouncement={handleApproveAnnouncement}
          onRejectAnnouncement={handleRejectAnnouncement}
          onApproveHandyman={handleApproveHandyman}
          onRejectHandyman={handleRejectHandyman}
          onBack={() => setCurrentPage('home')}
        />
      )}

      {currentPage === 'home' && (
        <HomePage 
          handymen={handymen}
          announcements={announcements}
          onSelectHandyman={handleSelectHandyman}
          userRole={userRole}
        />
      )}

      {currentPage === 'profile' && selectedHandymanId && (
        <ProfilePage 
          handyman={handymen.find(h => h.id === selectedHandymanId)!}
          onChat={handleChat} 
          onBack={handleBackToHome}
          onAddReview={handleAddReview}
        />
      )}

      {currentPage === 'chat' && selectedHandymanId && (
        <ChatPage 
          handymanId={selectedHandymanId} 
          onBack={handleBackToProfile} 
        />
      )}

      {currentPage === 'profile-user' && userRole === 'user' && (
        currentUserProfile ? (
          isEditingUserProfile ? (
            <UserProfileEditor
              profile={currentUserProfile}
              onSave={async (updatedProfile: UserProfile) => {
                try {
                  await api.updateProfile({
                    name: updatedProfile.name,
                    phone: updatedProfile.phone,
                    location: updatedProfile.location
                  });
                  setCurrentUserProfile(updatedProfile);
                  setIsEditingUserProfile(false);
                } catch (e) {
                  alert(e instanceof Error ? e.message : 'Error updating profile');
                }
              }}
              onBack={() => setIsEditingUserProfile(false)}
            />
          ) : (
            <UserProfilePage 
              announcements={announcements}
              onAddAnnouncement={handleAddAnnouncement}
              onBack={handleBackToHome}
              userProfile={currentUserProfile}
              onEditProfile={() => setIsEditingUserProfile(true)}
            />
          )
        ) : (
          showCreateUserProfile ? (
            <UserProfileEditor
              profile={{
                name: '',
                email: '',
                phone: '',
                location: '',
                memberSince: new Date().toLocaleDateString(),
              }}
              onSave={async (newProfile: UserProfile) => {
                try {
                  await api.updateProfile({
                    name: newProfile.name,
                    phone: newProfile.phone,
                    location: newProfile.location
                  });
                  setCurrentUserProfile(newProfile);
                  setShowCreateUserProfile(false);
                  setIsEditingUserProfile(false);
                  setCurrentPage('profile-user');
                } catch (e) {
                  alert(e instanceof Error ? e.message : 'Error creating profile');
                }
              }}
              onBack={() => setShowCreateUserProfile(false)}
            />
          ) : (
            <div className="max-w-2xl mx-auto py-8 text-center text-lg text-gray-500">
              Профиль пользователя не создан.<br />
              <button className="mt-6 px-6 py-3 bg-gray-800 text-white rounded-md font-bold hover:bg-gray-700 transition-all" onClick={() => setShowCreateUserProfile(true)}>
                Создать профиль пользователя
              </button>
            </div>
          )
        )
      )}
      {currentPage === 'profile-user' && userRole === 'handyman' && (
        currentOwnHandymanProfile ? (
          isEditingOwnHandymanProfile ? (
            <HandymanProfileEditor
              handyman={currentOwnHandymanProfile}
              onSave={async (updatedHandyman) => {
                try {
                  await api.updateProfile({
                    description: updatedHandyman.description,
                    category_names: [updatedHandyman.skill],
                    // add other fields if api supports
                  });
                  // re-fetch to make sure state is consistent
                  const p = await api.getProfile();
                  setCurrentUserProfile({
                    name: p.name || '',
                    email: p.email || '',
                    phone: p.phone || '',
                    location: p.worker?.location || '',
                    memberSince: new Date().toLocaleDateString(),
                    worker: p.worker,
                  });
                  
                  if (handymen.some(h => h.id === updatedHandyman.id)) {
                    setHandymen(prev => prev.map(h => h.id === updatedHandyman.id ? updatedHandyman : h));
                  }
                  if (pendingHandymen.some(h => h.id === updatedHandyman.id)) {
                    setPendingHandymen(prev => prev.map(h => h.id === updatedHandyman.id ? updatedHandyman : h));
                  }
                  setIsEditingOwnHandymanProfile(false);
                } catch (e) {
                  alert(e instanceof Error ? e.message : 'Error updating handyman profile');
                }
              }}
              onBack={() => setIsEditingOwnHandymanProfile(false)}
            />
          ) : (
            <ProfilePage
              handyman={currentOwnHandymanProfile}
              onBack={handleBackToHome}
              onAddReview={() => {}}
              isOwnProfile
              onEditProfile={() => setIsEditingOwnHandymanProfile(true)}
            />
          )
        ) : (
          showCreateHandyman ? (
            <HandymanProfileEditor 
              handyman={{
                id: Date.now().toString(),
                name: '',
                skill: 'Сантехника',
                rating: 0,
                reviews: [],
                hourlyRate: 0,
                description: '',
              }}
              onSave={newHandyman => {
                setPendingHandymen(prev => [newHandyman, ...prev]);
                setCurrentHandymanId(newHandyman.id);
                setShowCreateHandyman(false);
                setIsEditingOwnHandymanProfile(false);
                setCurrentPage('profile-user');
              }}
              onBack={() => setShowCreateHandyman(false)}
            />
          ) : (
            <div className="max-w-2xl mx-auto py-8 text-center text-lg text-gray-500">
              Профиль мастера не найден или ожидает модерацию.<br />
              <button className="mt-6 px-6 py-3 bg-gray-800 text-white rounded-md font-bold hover:bg-gray-700 transition-all" onClick={() => setShowCreateHandyman(true)}>
                Создать профиль мастера
              </button>
            </div>
          )
        )
      )}
    </Layout>
  );
}


