import React, { useState } from 'react';
import { Star, MessageSquare, ArrowLeft, Clock, Shield, Award, Send, Calendar } from 'lucide-react';
import { Handyman } from '../data';
import { ScheduleModal } from './ScheduleModal';
import { api } from '../api';

interface ProfilePageProps {
  handyman: Handyman;
  onChat?: (id: string) => void;
  onBack: () => void;
  onAddReview: (id: string, review: any) => void;
  isOwnProfile?: boolean;
  onEditProfile?: () => void;
  currentUserProfile?: any;
}

export function ProfilePage({ handyman, onChat, onBack, onAddReview, isOwnProfile = false, onEditProfile, currentUserProfile }: ProfilePageProps) {
  const [reviewText, setReviewText] = useState('');
  const [rating, setRating] = useState(5);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isScheduleOpen, setIsScheduleOpen] = useState(false);

  const handleReviewSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!reviewText.trim()) return;

    setIsSubmitting(true);
    // Frontend guard: prevent multiple reviews from same user
    const currentUserName = currentUserProfile?.name;
    if (currentUserName && handyman.reviews.some((r: any) => r.author === currentUserName)) {
      alert('Вы уже оставили отзыв для этого специалиста.');
      setIsSubmitting(false);
      return;
    }

    try {
      const res = await api.createReview({
        worker_id: Number(handyman.id),
        rating: rating,
        text: reviewText
      });
      
      const newReview = {
        id: String(res.id),
        author: res.author || 'Я',
        rating: res.rating,
        text: res.text,
        date: new Date(res.date).toLocaleDateString()
      };
      
      onAddReview(handyman.id, newReview);
      setReviewText('');
      setRating(5);
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Error creating review');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto pb-12">
      <button 
        onClick={onBack}
        className="mb-6 flex items-center text-gray-500 hover:text-gray-900 font-medium transition-colors"
      >
        <ArrowLeft className="w-4 h-4 mr-2" />
        Назад к списку
      </button>

      <div className="bg-white border-2 border-gray-200 rounded-lg overflow-hidden shadow-sm">
        {/* Header Section */}
        <div className="p-8 border-b-2 border-gray-100 flex flex-col md:flex-row gap-8 items-start">
          <div className="w-32 h-32 bg-gray-200 rounded-full flex-shrink-0 flex items-center justify-center border-4 border-white shadow-lg">
             <span className="text-4xl font-bold text-gray-400">{handyman.name.charAt(0)}</span>
          </div>
          
          <div className="flex-1">
            <div className="flex justify-between items-start">
              <div>
                <h1 className="text-3xl font-bold text-gray-900 mb-2">{handyman.name}</h1>
                <p className="text-lg text-gray-600 mb-4">{handyman.skill} Specialist</p>
              </div>
              <div className="flex flex-col items-end">
                <div className="text-2xl font-bold text-gray-900">{handyman.hourlyRate}<span className="text-sm font-normal text-gray-500">руб./час</span></div>
                <div className="flex items-center mt-1 text-sm font-bold">
                  <Star className="w-4 h-4 fill-yellow-400 text-yellow-400 mr-1" />
                  {handyman.rating} ({handyman.reviews.length} reviews)
                </div>
              </div>
            </div>

            <div className="flex gap-4 mt-2">
              {isOwnProfile ? (
                <button
                  onClick={onEditProfile}
                  className="bg-gray-800 text-white px-6 py-2.5 rounded-md font-bold hover:bg-gray-700 flex items-center gap-2 shadow-[4px_4px_0px_0px_rgba(0,0,0,0.1)] active:shadow-none active:translate-x-[2px] active:translate-y-[2px] transition-all"
                >
                  Редактировать профиль
                </button>
              ) : (
                <button 
                  onClick={() => onChat?.(handyman.id)}
                  className="bg-gray-800 text-white px-6 py-2.5 rounded-md font-bold hover:bg-gray-700 flex items-center gap-2 shadow-[4px_4px_0px_0px_rgba(0,0,0,0.1)] active:shadow-none active:translate-x-[2px] active:translate-y-[2px] transition-all"
                >
                  <MessageSquare className="w-4 h-4" />
                  Связаться с {handyman.name.split(' ')[0]}
                </button>
              )}
            </div>
          </div>
        </div>

        <div className="grid md:grid-cols-3 divide-y md:divide-y-0 md:divide-x divide-gray-100 border-b-2 border-gray-100">
           <div className="p-6 text-center">
             <Shield className="w-8 h-8 mx-auto mb-2 text-gray-400" />
             <h3 className="font-bold text-sm text-gray-900">Проверенный специалист</h3>
             <p className="text-xs text-gray-500 mt-1">Проверка фона и одобрение</p>
           </div>
           <div className="p-6 text-center">
             <Award className="w-8 h-8 mx-auto mb-2 text-gray-400" />
             <h3 className="font-bold text-sm text-gray-900">Гарантия удовлетворения</h3>
             <p className="text-xs text-gray-500 mt-1">Возврат денег при неудовлетворении</p>
           </div>
           <div className="p-6 text-center">
             <Clock className="w-8 h-8 mx-auto mb-2 text-gray-400" />
             <h3 className="font-bold text-sm text-gray-900">Быстрый ответ</h3>
             <p className="text-xs text-gray-500 mt-1">Отвечает в течение 2 часов</p>
           </div>
        </div>

        {/* Description & Reviews */}
        <div className="p-8 grid md:grid-cols-3 gap-8">
          <div className="md:col-span-2 space-y-12">
            <section>
              <h3 className="text-xl font-bold text-gray-900 mb-4">Обо мне</h3>
              <p className="text-gray-600 leading-relaxed">
                {handyman.description}
              </p>
            </section>

            <section>
              <div className="flex justify-between items-center mb-6">
                 <h3 className="text-xl font-bold text-gray-900">Отзывы клиентов</h3>
                 <span className="text-sm font-medium text-gray-500">{handyman.reviews.length} всего</span>
              </div>
              
              <div className="space-y-6 mb-8">
                {handyman.reviews.length > 0 ? (
                  handyman.reviews.map(review => (
                    <div key={review.id} className="bg-gray-50 p-4 rounded-lg border border-gray-100">
                      <div className="flex justify-between items-center mb-2">
                        <span className="font-bold text-gray-900">{review.author}</span>
                        <span className="text-xs text-gray-500">{review.date}</span>
                      </div>
                      <div className="flex items-center mb-2">
                         {[...Array(5)].map((_, i) => (
                           <Star 
                             key={i} 
                             className={`w-3 h-3 ${i < review.rating ? 'fill-gray-900 text-gray-900' : 'text-gray-300'}`} 
                           />
                         ))}
                      </div>
                      <p className="text-gray-600 text-sm">{review.text}</p>
                    </div>
                  ))
                ) : (
                  <p className="text-gray-500 italic">Отзывов пока нет.</p>
                )}
              </div>

              {!isOwnProfile && (
                <div className="bg-gray-50 p-6 rounded-lg border-2 border-gray-200">
                  <h4 className="font-bold text-gray-900 mb-4">Оставить отзыв</h4>
                  <form onSubmit={handleReviewSubmit}>
                    <div className="mb-4">
                      <label className="block text-sm font-bold text-gray-700 mb-2">Рейтинг</label>
                      <div className="flex gap-1">
                        {[1, 2, 3, 4, 5].map((star) => (
                          <button
                            key={star}
                            type="button"
                            onClick={() => setRating(star)}
                            className="focus:outline-none transition-transform active:scale-95"
                          >
                            <Star 
                              className={`w-6 h-6 ${
                                star <= rating 
                                  ? 'fill-yellow-400 text-yellow-400' 
                                  : 'text-gray-300 hover:text-gray-400'
                              }`} 
                            />
                          </button>
                        ))}
                      </div>
                    </div>
                    
                    <div className="mb-4">
                      <label className="block text-sm font-bold text-gray-700 mb-2">Отзыв</label>
                      <textarea 
                        value={reviewText}
                        onChange={(e) => setReviewText(e.target.value)}
                        className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none resize-none bg-white"
                        rows={3}
                        placeholder="Поделитесь своим опытом..."
                        minLength={10}
                        maxLength={1000}
                        required
                      />
                    </div>
                    
                    <button 
                      type="submit"
                      disabled={isSubmitting || !reviewText.trim()}
                      className="bg-gray-900 text-white px-6 py-2 rounded font-bold hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
                    >
                      {isSubmitting ? 'Публикация...' : 'Опубликовать отзыв'}
                      {!isSubmitting && <Send className="w-4 h-4" />}
                    </button>
                  </form>
                </div>
              )}
            </section>
          </div>
          
          <div className="bg-gray-50 p-6 rounded-lg h-fit border border-gray-200">
            <h4 className="font-bold text-gray-900 mb-4">Доступность</h4>
            <div className="space-y-3 text-sm">
               <div className="flex justify-between">
                 <span className="text-gray-600">Пн - Пт</span>
                 <span className="font-medium">9:00 - 18:00</span>
               </div>
               <div className="flex justify-between">
                 <span className="text-gray-600">Суббота</span>
                 <span className="font-medium">10:00 - 16:00</span>
               </div>
               <div className="flex justify-between">
                 <span className="text-gray-600">Воскресенье</span>
                 <span className="font-medium text-red-500">Закрыто</span>
               </div>
            </div>
            <button 
              onClick={() => setIsScheduleOpen(true)}
              className="w-full mt-6 border-2 border-gray-800 text-gray-800 font-bold py-2 rounded text-sm hover:bg-gray-200 flex items-center justify-center gap-2 transition-colors"
            >
              <Calendar className="w-4 h-4" />
              Проверить расписание
            </button>
          </div>
        </div>
      </div>

      <ScheduleModal 
        isOpen={isScheduleOpen}
        onClose={() => setIsScheduleOpen(false)}
        handymanName={handyman.name}
      />
    </div>
  );
}
