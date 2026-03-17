import React, { useState } from 'react';
import { X, ChevronLeft, ChevronRight, Clock } from 'lucide-react';

interface ScheduleModalProps {
  isOpen: boolean;
  onClose: () => void;
  handymanName: string;
}

export function ScheduleModal({ isOpen, onClose, handymanName }: ScheduleModalProps) {
  if (!isOpen) return null;

  const [selectedDate, setSelectedDate] = useState<number | null>(null);
  const [currentMonth, setCurrentMonth] = useState('Ноябрь 2023');

  // Mock calendar days
  const days = Array.from({ length: 30 }, (_, i) => i + 1);
  const slots = [
    '09:00', '10:00', '11:00', '13:00', '14:00', '16:00'
  ];

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md overflow-hidden animate-in fade-in zoom-in duration-200">
        <div className="flex justify-between items-center p-4 border-b border-gray-100">
          <h3 className="font-bold text-lg text-gray-900">Проверить доступность</h3>
          <button onClick={onClose} className="p-1 hover:bg-gray-100 rounded-full transition-colors">
            <X className="w-5 h-5 text-gray-500" />
          </button>
        </div>
        
        <div className="p-6">
          <p className="text-sm text-gray-500 mb-4">
            Просмотр расписания для <span className="font-bold text-gray-900">{handymanName}</span>
          </p>

          {/* Month Navigation */}
          <div className="flex justify-between items-center mb-4">
            <button className="p-1 hover:bg-gray-100 rounded">
              <ChevronLeft className="w-5 h-5 text-gray-600" />
            </button>
            <span className="font-bold text-gray-900">{currentMonth}</span>
            <button className="p-1 hover:bg-gray-100 rounded">
              <ChevronRight className="w-5 h-5 text-gray-600" />
            </button>
          </div>

          {/* Calendar Grid */}
          <div className="grid grid-cols-7 gap-1 mb-6">
            {['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'].map(day => (
              <div key={day} className="text-center text-xs font-bold text-gray-400 py-1">
                {day}
              </div>
            ))}
            {days.map(day => (
              <button
                key={day}
                onClick={() => setSelectedDate(day)}
                className={`aspect-square flex items-center justify-center text-sm rounded-md transition-colors ${
                  selectedDate === day
                    ? 'bg-gray-900 text-white font-bold'
                    : 'hover:bg-gray-100 text-gray-700'
                }`}
              >
                {day}
              </button>
            ))}
          </div>

          {/* Time Slots */}
          <div className="border-t border-gray-100 pt-4">
            <h4 className="text-sm font-bold text-gray-900 mb-3 flex items-center">
              <Clock className="w-4 h-4 mr-2" />
              Доступные слоты {selectedDate ? `на ${selectedDate} ноября` : ''}
            </h4>
            <div className="grid grid-cols-3 gap-2">
              {selectedDate ? (
                slots.map((slot, index) => (
                  <button 
                    key={index}
                    className="px-2 py-2 text-xs font-medium border border-gray-200 rounded hover:border-gray-900 hover:bg-gray-50 transition-all text-center"
                  >
                    {slot}
                  </button>
                ))
              ) : (
                <div className="col-span-3 text-center text-gray-400 text-sm py-4 italic">
                  Выберите дату, чтобы увидеть доступное время
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="p-4 bg-gray-50 border-t border-gray-100 text-center">
          <button 
            onClick={onClose}
            className="text-sm font-bold text-gray-600 hover:text-gray-900"
          >
            Закрыть
          </button>
        </div>
      </div>
    </div>
  );
}
