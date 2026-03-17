import React, { useState } from 'react';
import { User, Mail, Phone, MapPin, Edit, Briefcase, Plus, Calendar, Clock, ArrowRight } from 'lucide-react';
import { CreateAnnouncementModal } from './CreateAnnouncementModal';
import { Announcement } from '../data';

interface UserProfilePageProps {
  onBack: () => void;
  onEditProfile?: () => void;
  announcements: Announcement[];
  onAddAnnouncement: (announcement: Announcement) => void;
  userProfile: {
    name: string;
    email: string;
    phone: string;
    location: string;
    memberSince: string;
  };
}

export function UserProfilePage({ onBack, onEditProfile, announcements, onAddAnnouncement, userProfile }: UserProfilePageProps) {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [activeTab, setActiveTab] = useState<'active' | 'history'>('active');
  const [myAnnouncements, setMyAnnouncements] = useState<Announcement[]>([]);
  const isHandyman = false; // Здесь можно добавить проверку роли, если потребуется

  // Load user specific ads from the backend
  React.useEffect(() => {
    import('../api').then(({ api }) => {
      api.getMyAds().then((res) => {
        if (res && Array.isArray(res)) {
           const mappedAnnouncements: Announcement[] = res.map((ad: any) => ({
            id: ad.id,
            title: ad.title,
            category: ad.category?.name || 'Без категории',
            status: ad.status || 'На модерации',
            date: ad.created_at ? new Date(ad.created_at).toLocaleDateString() : '',
            budget: `${ad.price} руб.`,
            handyman: ad.user?.name || null,
            location: ad.location,
            description: '', 
          }));
          setMyAnnouncements(mappedAnnouncements);
        }
      }).catch(console.error);
    });
  }, []);

  const handleCreateAnnouncement = (data: any) => {
    // data уже содержит ответ API (createAd), если запрос успешен
    const newAnnouncement: Announcement = {
      id: data.id || Date.now(),
      title: data.title,
      category: data.category?.name || 'Без категории',
      status: 'На модерации',
      date: data.created_at ? new Date(data.created_at).toLocaleDateString() : new Date().toLocaleDateString(),
      budget: `${data.price ?? data.budget ?? 0} руб.`,
      handyman: null,
      location: data.location,
      description: '', // серверная модель объявления не хранит отдельного описания
    };
    onAddAnnouncement(newAnnouncement);
    setMyAnnouncements(prev => [newAnnouncement, ...prev]);
    setIsModalOpen(false);
  };

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      {/* Back Button */}
      <button 
        onClick={onBack}
        className="mb-6 flex items-center text-gray-500 hover:text-gray-900 font-medium transition-colors"
      >
        <ArrowRight className="w-4 h-4 mr-2 rotate-180" />
        Назад к панели управления
      </button>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        {/* Left Column: User Details Card */}
        <div className="lg:col-span-1">
          <div className="bg-white border-2 border-gray-200 rounded-lg overflow-hidden shadow-sm sticky top-24">
            {/* Dark Header Banner */}
            <div className="h-32 bg-[#1F2937] relative">
              <button className="absolute top-4 right-4 p-2 bg-white/10 hover:bg-white/20 rounded-md text-white transition-colors border border-white/20">
                <Edit className="w-4 h-4" />
              </button>
            </div>
            
            <div className="px-6 pb-6 relative">
              {/* Avatar */}
              <div className="absolute -top-12 left-6 w-24 h-24 bg-white p-1 rounded-lg shadow-sm border border-gray-100">
                <div className="w-full h-full bg-gray-200 rounded flex items-center justify-center text-gray-500">
                  <User className="w-10 h-10" />
                </div>
              </div>
              
              {/* Content */}
              <div className="pt-16 space-y-6">
                <div>
                  <h2 className="text-2xl font-bold text-gray-900 leading-tight">{userProfile.name}</h2>
                  <p className="text-gray-500 text-sm mt-1">Участник с {userProfile.memberSince}</p>
                </div>

                <div className="space-y-4 pt-2">
                  <div className="flex items-center text-gray-600 text-sm group cursor-pointer hover:text-gray-900">
                    <Mail className="w-4 h-4 mr-3 text-gray-400 group-hover:text-gray-600" />
                    {userProfile.email}
                  </div>
                  <div className="flex items-center text-gray-600 text-sm group cursor-pointer hover:text-gray-900">
                    <Phone className="w-4 h-4 mr-3 text-gray-400 group-hover:text-gray-600" />
                    {userProfile.phone}
                  </div>
                  <div className="flex items-center text-gray-600 text-sm group cursor-pointer hover:text-gray-900">
                    <MapPin className="w-4 h-4 mr-3 text-gray-400 group-hover:text-gray-600" />
                    {userProfile.location}
                  </div>
                </div>

                <button onClick={onEditProfile} className="w-full bg-white border border-gray-300 text-gray-700 font-bold py-2.5 rounded hover:bg-gray-50 transition-colors text-sm shadow-sm active:translate-y-0.5">
                  Редактировать профиль
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column: Announcements & History */}
        <div className="lg:col-span-2 space-y-6">
          
          {/* Action Header (только для пользователя) */}
          {!isHandyman && (
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 bg-gray-50 p-6 rounded-lg border border-gray-200">
              <div>
                <h3 className="text-lg font-bold text-gray-900">Нужно что-то сделать?</h3>
                <p className="text-gray-500 text-sm">Опубликуйте новую работу и получайте предложения от профессионалов.</p>
              </div>
              <button 
                onClick={() => setIsModalOpen(true)}
                className="bg-gray-900 text-white px-5 py-2.5 rounded font-bold hover:bg-gray-800 flex items-center gap-2 shadow-[4px_4px_0px_0px_rgba(0,0,0,0.1)] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none transition-all"
              >
                <Plus className="w-5 h-5" />
                Создать объявление
              </button>
            </div>
          )}

          {/* Tabs */}
          <div className="bg-white border-2 border-gray-200 rounded-lg overflow-hidden min-h-[400px]">
            <div className="flex border-b-2 border-gray-100">
              <button 
                onClick={() => setActiveTab('active')}
                className={`flex-1 py-4 text-center font-bold text-sm border-b-2 transition-colors ${
                  activeTab === 'active' 
                    ? 'border-gray-900 text-gray-900 bg-gray-50' 
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                }`}
              >
                Мои объявления
              </button>
              <button 
                onClick={() => setActiveTab('history')}
                className={`flex-1 py-4 text-center font-bold text-sm border-b-2 transition-colors ${
                  activeTab === 'history' 
                    ? 'border-gray-900 text-gray-900 bg-gray-50' 
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                }`}
              >
                История
              </button>
            </div>

            {/* List Content */}
            <div className="p-0">
              {myAnnouncements.length === 0 ? (
                <div className="p-12 text-center text-gray-400">
                  <Briefcase className="w-12 h-12 mx-auto mb-3 opacity-20" />
                  <p>Нет активных объявлений.</p>
                </div>
              ) : (
                <div className="divide-y divide-gray-100">
                  {myAnnouncements.map((announcement) => (
                    <div key={announcement.id} className="p-6 hover:bg-gray-50 transition-colors group">
                      <div className="flex justify-between items-start mb-2">
                        <div>
                           <h4 className="font-bold text-gray-900 text-lg group-hover:text-blue-600 transition-colors">
                             {announcement.title}
                           </h4>
                           <span className="text-xs font-bold bg-gray-200 text-gray-700 px-2 py-0.5 rounded mt-1 inline-block">
                             {announcement.category}
                           </span>
                        </div>
                        <div className={`px-3 py-1 rounded text-xs font-bold uppercase tracking-wider ${
                          announcement.status === 'Open' ? 'bg-green-100 text-green-700' : 'bg-blue-100 text-blue-700'
                        }`}>
                          {announcement.status === 'Open' ? 'Открыто' : announcement.status}
                        </div>
                      </div>
                      
                      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-4 text-sm text-gray-500">
                        <div className="flex items-center">
                          <Calendar className="w-4 h-4 mr-2" />
                          {announcement.date}
                        </div>
                        <div className="flex items-center font-mono">
                          {announcement.budget}
                        </div>
                        <div className="col-span-2 flex items-center sm:justify-end">
                           {announcement.handyman ? (
                             <span className="flex items-center text-gray-900 font-medium">
                               <User className="w-4 h-4 mr-2" />
                               Назначено: {announcement.handyman}
                             </span>
                           ) : (
                             <span className="text-gray-400 italic">Мастер не назначен</span>
                           )}
                        </div>
                      </div>

                      <div className="mt-4 flex gap-2 justify-end opacity-0 group-hover:opacity-100 transition-opacity">
                        <button className="text-xs font-bold text-gray-500 hover:text-gray-900 underline">Просмотреть детали</button>
                        <span className="text-gray-300">|</span>
                        <button className="text-xs font-bold text-red-500 hover:text-red-700 underline">Отменить</button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Modal */}
      <CreateAnnouncementModal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateAnnouncement}
      />
    </div>
  );
}
