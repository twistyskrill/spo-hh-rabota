import React, { useEffect, useMemo, useState } from 'react';
import { X, ChevronLeft, ChevronRight, Clock, Check } from 'lucide-react';
import { api, AvailabilitySlot } from '../api';

interface ScheduleModalProps {
  isOpen: boolean;
  onClose: () => void;
  handymanName: string;
  handymanId: number;
  isOwnProfile?: boolean;
}

const WEEK_DAYS = ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'];
const MONTHS_RU = ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'];
const TIME_OPTIONS = ['08:00', '09:00', '10:00', '11:00', '12:00', '13:00', '14:00', '15:00', '16:00', '17:00', '18:00', '19:00', '20:00'];

const toDayKey = (date: Date): string => {
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, '0');
  const d = String(date.getDate()).padStart(2, '0');
  return `${y}-${m}-${d}`;
};

const toTimeKey = (date: Date): string => {
  const h = String(date.getHours()).padStart(2, '0');
  const m = String(date.getMinutes()).padStart(2, '0');
  return `${h}:${m}`;
};

const toSlotKey = (iso: string): string => {
  const date = new Date(iso);
  return `${toDayKey(date)}|${toTimeKey(date)}`;
};

const keyToIso = (key: string): string => {
  const [day, time] = key.split('|');
  const [y, m, d] = day.split('-').map(Number);
  const [h, mm] = time.split(':').map(Number);
  return new Date(y, m - 1, d, h, mm, 0, 0).toISOString();
};

const buildMonthGrid = (monthDate: Date): Array<Date | null> => {
  const year = monthDate.getFullYear();
  const month = monthDate.getMonth();
  const firstDay = new Date(year, month, 1);
  const daysInMonth = new Date(year, month + 1, 0).getDate();
  const leading = firstDay.getDay();

  const cells: Array<Date | null> = [];
  for (let i = 0; i < leading; i += 1) {
    cells.push(null);
  }
  for (let day = 1; day <= daysInMonth; day += 1) {
    cells.push(new Date(year, month, day));
  }
  return cells;
};

export function ScheduleModal({ isOpen, onClose, handymanName, handymanId, isOwnProfile = false }: ScheduleModalProps) {
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  const [currentMonth, setCurrentMonth] = useState(new Date());
  const [slots, setSlots] = useState<AvailabilitySlot[]>([]);
  const [editableSlotKeys, setEditableSlotKeys] = useState<Set<string>>(new Set());
  const [selectedSlotKey, setSelectedSlotKey] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const monthLabel = `${MONTHS_RU[currentMonth.getMonth()]} ${currentMonth.getFullYear()}`;
  const monthGrid = useMemo(() => buildMonthGrid(currentMonth), [currentMonth]);

  const rangeFrom = useMemo(
    () => new Date(currentMonth.getFullYear(), currentMonth.getMonth(), 1, 0, 0, 0, 0).toISOString(),
    [currentMonth]
  );
  const rangeTo = useMemo(
    () => new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1, 0, 23, 59, 59, 999).toISOString(),
    [currentMonth]
  );

  useEffect(() => {
    if (!isOpen) return;

    const fetchSlots = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = isOwnProfile
          ? await api.getMyAvailability(rangeFrom, rangeTo)
          : await api.getHandymanAvailability(handymanId, rangeFrom, rangeTo);
        const fetchedSlots = response?.slots || [];
        setSlots(fetchedSlots);

        if (isOwnProfile) {
          const freeKeys = fetchedSlots.filter(s => !s.is_booked).map(s => toSlotKey(s.start_at));
          setEditableSlotKeys(new Set(freeKeys));
        }
      } catch (e) {
        setError(e instanceof Error ? e.message : 'Ошибка загрузки расписания');
      } finally {
        setLoading(false);
      }
    };

    fetchSlots();
  }, [isOpen, isOwnProfile, handymanId, rangeFrom, rangeTo]);

  const bookedSlotKeys = useMemo(
    () => new Set(slots.filter(s => s.is_booked).map(s => toSlotKey(s.start_at))),
    [slots]
  );

  const availableDayKeys = useMemo(() => {
    if (isOwnProfile) {
      const dayKeys = new Set<string>();
      editableSlotKeys.forEach(key => dayKeys.add(key.split('|')[0]));
      bookedSlotKeys.forEach(key => dayKeys.add(key.split('|')[0]));
      return dayKeys;
    }

    const dayKeys = new Set<string>();
    slots.forEach(slot => {
      if (!slot.is_booked) {
        dayKeys.add(toDayKey(new Date(slot.start_at)));
      }
    });
    return dayKeys;
  }, [isOwnProfile, editableSlotKeys, bookedSlotKeys, slots]);

  const selectedDayKey = selectedDate ? toDayKey(selectedDate) : null;

  const userDaySlots = useMemo(() => {
    if (!selectedDayKey || isOwnProfile) return [];

    return slots
      .filter(slot => !slot.is_booked && toDayKey(new Date(slot.start_at)) === selectedDayKey)
      .sort((a, b) => new Date(a.start_at).getTime() - new Date(b.start_at).getTime());
  }, [slots, selectedDayKey, isOwnProfile]);

  const handleToggleMasterSlot = (timeValue: string) => {
    if (!selectedDate) return;
    const slotKey = `${toDayKey(selectedDate)}|${timeValue}`;
    if (bookedSlotKeys.has(slotKey)) return;

    setEditableSlotKeys(prev => {
      const next = new Set(prev);
      if (next.has(slotKey)) {
        next.delete(slotKey);
      } else {
        next.add(slotKey);
      }
      return next;
    });
  };

  const handleSaveMasterAvailability = async () => {
    setSaving(true);
    setError(null);
    try {
      const payload = Array.from(editableSlotKeys).map(keyToIso);
      await api.updateMyAvailability(payload);
      alert('Расписание обновлено');
      onClose();
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Не удалось сохранить расписание');
    } finally {
      setSaving(false);
    }
  };

  const handleBookSlot = async () => {
    if (!selectedSlotKey) return;

    setSaving(true);
    setError(null);
    try {
      await api.bookHandymanSlot(handymanId, keyToIso(selectedSlotKey));
      alert('Бронирование успешно создано');
      setSlots(prev => prev.filter(s => toSlotKey(s.start_at) !== selectedSlotKey));
      setSelectedSlotKey(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Не удалось забронировать время');
    } finally {
      setSaving(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-lg overflow-hidden animate-in fade-in zoom-in duration-200">
        <div className="flex justify-between items-center p-4 border-b border-gray-100">
          <h3 className="font-bold text-lg text-gray-900">{isOwnProfile ? 'Мой календарь' : 'Проверить доступность'}</h3>
          <button onClick={onClose} className="p-1 hover:bg-gray-100 rounded-full transition-colors">
            <X className="w-5 h-5 text-gray-500" />
          </button>
        </div>
        
        <div className="p-6">
          <p className="text-sm text-gray-500 mb-4">
            {isOwnProfile
              ? 'Выберите даты и часы, когда вы готовы начать работу.'
              : <>Просмотр свободных слотов для <span className="font-bold text-gray-900">{handymanName}</span></>}
          </p>

          <div className="flex justify-between items-center mb-4">
            <button
              className="p-1 hover:bg-gray-100 rounded"
              onClick={() => setCurrentMonth(prev => new Date(prev.getFullYear(), prev.getMonth() - 1, 1))}
            >
              <ChevronLeft className="w-5 h-5 text-gray-600" />
            </button>
            <span className="font-bold text-gray-900">{monthLabel}</span>
            <button
              className="p-1 hover:bg-gray-100 rounded"
              onClick={() => setCurrentMonth(prev => new Date(prev.getFullYear(), prev.getMonth() + 1, 1))}
            >
              <ChevronRight className="w-5 h-5 text-gray-600" />
            </button>
          </div>

          <div className="grid grid-cols-7 gap-1 mb-6">
            {WEEK_DAYS.map(day => (
              <div key={day} className="text-center text-xs font-bold text-gray-400 py-1">
                {day}
              </div>
            ))}
            {monthGrid.map((dateCell, idx) => {
              if (!dateCell) {
                return <div key={`empty-${idx}`} className="aspect-square" />;
              }

              const dayKey = toDayKey(dateCell);
              const isSelected = selectedDayKey === dayKey;
              const hasAvailability = availableDayKeys.has(dayKey);
              const isDisabledForUser = !isOwnProfile && !hasAvailability;

              return (
                <button
                  key={dayKey}
                  onClick={() => {
                    setSelectedDate(dateCell);
                    setSelectedSlotKey(null);
                  }}
                  disabled={isDisabledForUser}
                  className={`aspect-square flex items-center justify-center text-sm rounded-md transition-colors ${
                    isSelected
                      ? 'bg-gray-900 text-white font-bold'
                      : isDisabledForUser
                        ? 'text-gray-300 cursor-not-allowed'
                        : hasAvailability
                          ? 'bg-gray-100 text-gray-900 hover:bg-gray-200 font-medium'
                          : 'hover:bg-gray-100 text-gray-700'
                  }`}
                >
                  {dateCell.getDate()}
                </button>
              );
            })}
          </div>

          <div className="border-t border-gray-100 pt-4">
            <h4 className="text-sm font-bold text-gray-900 mb-3 flex items-center">
              <Clock className="w-4 h-4 mr-2" />
              {isOwnProfile ? 'Часы доступности' : 'Свободные слоты'}
            </h4>
            {loading ? (
              <div className="text-center text-gray-500 text-sm py-4">Загрузка...</div>
            ) : (
              <div className="grid grid-cols-3 gap-2">
                {!selectedDate ? (
                  <div className="col-span-3 text-center text-gray-400 text-sm py-4 italic">
                    Выберите дату, чтобы увидеть доступное время
                  </div>
                ) : isOwnProfile ? (
                  TIME_OPTIONS.map(timeValue => {
                    const slotKey = `${toDayKey(selectedDate)}|${timeValue}`;
                    const isBooked = bookedSlotKeys.has(slotKey);
                    const selected = editableSlotKeys.has(slotKey);
                    return (
                      <button
                        key={timeValue}
                        disabled={isBooked}
                        onClick={() => handleToggleMasterSlot(timeValue)}
                        className={`px-2 py-2 text-xs font-medium border rounded transition-all text-center ${
                          isBooked
                            ? 'bg-gray-100 text-gray-400 border-gray-200 cursor-not-allowed'
                            : selected
                              ? 'bg-gray-900 text-white border-gray-900'
                              : 'border-gray-200 hover:border-gray-900 hover:bg-gray-50'
                        }`}
                      >
                        {timeValue}{isBooked ? ' занято' : ''}
                      </button>
                    );
                  })
                ) : userDaySlots.length > 0 ? (
                  userDaySlots.map(slot => {
                    const slotKey = toSlotKey(slot.start_at);
                    const selected = selectedSlotKey === slotKey;
                    const timeLabel = toTimeKey(new Date(slot.start_at));
                    return (
                      <button
                        key={slot.id}
                        onClick={() => setSelectedSlotKey(slotKey)}
                        className={`px-2 py-2 text-xs font-medium border rounded transition-all text-center ${
                          selected
                            ? 'bg-gray-900 text-white border-gray-900'
                            : 'border-gray-200 hover:border-gray-900 hover:bg-gray-50'
                        }`}
                      >
                        {timeLabel}
                      </button>
                    );
                  })
                ) : (
                  <div className="col-span-3 text-center text-gray-400 text-sm py-4 italic">
                    На выбранную дату свободных слотов нет
                  </div>
                )}
                {error && (
                  <div className="col-span-3 text-center text-red-500 text-sm py-2">
                    {error}
                  </div>
                )}
              </div>
            )}
          </div>
        </div>

        <div className="p-4 bg-gray-50 border-t border-gray-100 flex items-center justify-between gap-3">
          <button onClick={onClose} className="text-sm font-bold text-gray-600 hover:text-gray-900">
            Отмена
          </button>

          {isOwnProfile ? (
            <button
              onClick={handleSaveMasterAvailability}
              disabled={saving || loading}
              className="bg-gray-900 text-white px-4 py-2 rounded text-sm font-bold hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed inline-flex items-center gap-2"
            >
              <Check className="w-4 h-4" />
              {saving ? 'Сохранение...' : 'Сохранить расписание'}
            </button>
          ) : (
            <button
              onClick={handleBookSlot}
              disabled={saving || loading || !selectedSlotKey}
              className="bg-gray-900 text-white px-4 py-2 rounded text-sm font-bold hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {saving ? 'Бронирование...' : 'Забронировать'}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
