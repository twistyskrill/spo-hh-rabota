import React, { useState, useEffect } from 'react';
import { Star, ArrowRight, Briefcase, MapPin, Calendar, DollarSign, User } from 'lucide-react';
import { Handyman, Announcement } from '../data';

interface HomePageProps {
  handymen: Handyman[];
  announcements: Announcement[];
  onSelectHandyman: (id: string) => void;
  userRole?: 'user' | 'handyman' | null;
}

export function HomePage({ handymen, announcements, onSelectHandyman, userRole = 'user' }: HomePageProps) {
  // Default view based on role
  const [viewMode, setViewMode] = useState<'pros' | 'jobs'>('pros');

  useEffect(() => {
    if (userRole === 'handyman') {
      setViewMode('jobs');
    } else {
      setViewMode('pros');
    }
  }, [userRole]);

  const [respondingId, setRespondingId] = useState<number | null>(null);
  const [respondedIds, setRespondedIds] = useState<Set<number>>(new Set());

  const handleRespond = async (jobId: number) => {
    setRespondingId(jobId);
    try {
      const api = (await import('../api')).api;
      await api.createResponse({ adId: jobId, message: "Откликнусь на вашу заявку", price: 0 });
      setRespondedIds(prev => new Set(prev).add(jobId));
    } catch (e) {
      console.error('Failed to respond to ad', e);
      alert('Ошибка при отклике');
    } finally {
      setRespondingId(null);
    }
  };

  return (
    <div className="space-y-8">
      {/* Header Section */}
      <div className="flex flex-col md:flex-row justify-between items-center border-b border-gray-200 pb-4 gap-4">
        <h2 className="text-2xl font-bold text-gray-900">
          {viewMode === 'pros' ? 'Найти специалистов' : 'Доска вакансий'}
        </h2>

        {/* Filters - Keep simple for now */}
        <div className="flex gap-2 w-full md:w-auto">
           <span className="text-sm font-medium text-gray-500 self-center hidden md:inline">Фильтр по:</span>
           <select className="text-sm border-gray-300 border rounded px-2 py-1 bg-white w-full md:w-auto">
             <option>Все категории</option>
             <option>Сантехника</option>
             <option>Электрика</option>
             <option>Столярные работы</option>
           </select>
        </div>
      </div>

      {viewMode === 'pros' ? (
        /* Handyman Grid (Visible to Users) */
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {handymen.map((handyman) => (
            <div 
              key={handyman.id} 
              className="group bg-white border-2 border-gray-200 rounded-lg overflow-hidden hover:border-gray-800 hover:shadow-[4px_4px_0px_0px_rgba(31,41,55,1)] transition-all duration-200 cursor-pointer"
              onClick={() => onSelectHandyman(handyman.id)}
            >
              <div className="h-48 bg-gray-100 flex items-center justify-center border-b-2 border-gray-100 group-hover:border-gray-800 transition-colors">
                <span className="text-gray-400 font-bold text-4xl select-none">
                  {handyman.name.charAt(0)}
                </span>
              </div>
              
              <div className="p-5">
                <div className="flex justify-between items-start mb-2">
                  <h3 className="font-bold text-lg text-gray-900 group-hover:text-blue-600 transition-colors">
                    {handyman.name}
                  </h3>
                  <div className="flex items-center bg-gray-100 px-2 py-1 rounded text-xs font-bold">
                    <Star className="w-3 h-3 fill-gray-900 text-gray-900 mr-1" />
                    {handyman.rating}
                  </div>
                </div>
                
                <div className="mb-4">
                  <span className="inline-block bg-gray-200 text-gray-700 text-xs px-2 py-1 rounded font-medium">
                    {handyman.skill}
                  </span>
                  <span className="ml-2 text-xs text-gray-500 font-mono">
                    {handyman.hourlyRate} руб./час
                  </span>
                </div>

                <button 
                  className="w-full flex items-center justify-center gap-2 border-2 border-gray-800 text-gray-800 font-bold py-2 rounded hover:bg-gray-800 hover:text-white transition-all text-sm"
                >
                  Открыть профиль
                  <ArrowRight className="w-4 h-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      ) : (
        /* Job Board List (Visible to Handymen) */
        <div className="space-y-4">
          {announcements.map((job) => (
            <div key={job.id} className="bg-white border-2 border-gray-200 rounded-lg p-6 hover:border-gray-800 transition-colors">
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-bold text-lg text-gray-900">{job.title}</h3>
                  <div className="flex items-center gap-2 mt-1">
                    <span className="bg-blue-100 text-blue-800 text-xs font-bold px-2 py-0.5 rounded">
                      {job.category}
                    </span>
                    <span className="text-xs text-gray-500 flex items-center">
                      <MapPin className="w-3 h-3 mr-1" />
                      {job.location || 'Удаленно'}
                    </span>
                  </div>
                </div>
                <div className="font-mono font-bold text-gray-900 bg-gray-50 px-3 py-1 rounded border border-gray-200">
                  {job.budget}
                </div>
              </div>
              
              <p className="text-gray-600 text-sm mt-3 line-clamp-2">
                {job.description || 'Описание не предоставлено.'}
              </p>
              
              <div className="flex items-center justify-between mt-4 pt-4 border-t border-gray-100">
                <div className="flex gap-4 text-xs text-gray-500 font-medium">
                   <span className="flex items-center">
                     <User className="w-3 h-3 mr-1" />
                     Опубликовано клиентом
                   </span>
                   <span className="flex items-center">
                     <Calendar className="w-3 h-3 mr-1" />
                     {job.date}
                   </span>
                </div>
                
                {respondedIds.has(Number(job.id)) ? (
                  <span className="text-sm font-bold text-green-600 flex items-center">
                    Вы откликнулись ✓
                  </span>
                ) : (
                  <button 
                    onClick={() => handleRespond(Number(job.id))}
                    disabled={respondingId === Number(job.id)}
                    className="text-sm font-bold text-gray-900 hover:underline flex items-center disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {respondingId === Number(job.id) ? 'Отправка...' : 'Откликнуться сейчас'} 
                    {respondingId !== Number(job.id) && <ArrowRight className="w-4 h-4 ml-1" />}
                  </button>
                )}
              </div>
            </div>
          ))}
          
          {announcements.length === 0 && (
            <div className="text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-200">
              <Briefcase className="w-12 h-12 text-gray-300 mx-auto mb-3" />
              <p className="text-gray-500 font-medium">Активных вакансий не найдено.</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
