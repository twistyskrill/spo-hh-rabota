import React, { useEffect, useState } from 'react';
import { X, Upload, MapPin, Briefcase } from 'lucide-react';
import { api } from '../api';
import { AIAssistButton } from './AIAssistButton';

interface CreateAnnouncementModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: any) => void;
}

export function CreateAnnouncementModal({ isOpen, onClose, onSubmit }: CreateAnnouncementModalProps) {
  const [formData, setFormData] = useState({
    title: '',
    categoryId: 0,
    priceUnitId: 0,
    budget: '',
    location: '',
    schedule: '',
    description: ''
  });

  const [categories, setCategories] = useState<{ id: number; name: string }[]>([]);
  const [priceUnits, setPriceUnits] = useState<{ id: number; name: string }[]>([]);

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!isOpen) return;

    const loadInfo = async () => {
      try {
        const [cats, units] = await Promise.all([
          api.getCategories(),
          api.getPriceUnits(),
        ]);
        setCategories(cats);
        setPriceUnits(units);
        setFormData(prev => ({
          ...prev,
          categoryId: cats[0]?.id ?? 0,
          priceUnitId: units[0]?.id ?? 0,
        }));
      } catch (e) {
        console.error(e);
        setError('Не удалось загрузить справочники категорий и единиц цены.');
      }
    };
    loadInfo();
  }, [isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');
    
    try {
      // simple client-side validation to avoid 400 от бэка
      const price = Number(formData.budget);
      if (
        !formData.title.trim() ||
        !formData.location.trim() ||
        !formData.categoryId ||
        !formData.priceUnitId ||
        !price ||
        price <= 0
      ) {
        setError('Заполните название, местоположение и бюджет (больше 0).');
        setIsLoading(false);
        return;
      }

      // Send the request using API (см. CreateAdRequest в бэкенде)
      const response = await api.createAd({
        title: formData.title,
        price,
        category_id: formData.categoryId,
        price_unit_id: formData.priceUnitId,
        location: formData.location,
        schedule: formData.schedule,
      });
      // Pass back either the API response or the original formData
      onSubmit({ ...formData, ...response });
      onClose();
    } catch (err) {
      console.error(err);
      setError('Не удалось создать объявление.');
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4 backdrop-blur-sm">
      <div className="bg-white w-full max-w-lg border-2 border-gray-900 shadow-[8px_8px_0px_0px_rgba(31,41,55,1)] flex flex-col max-h-[90vh]">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b-2 border-gray-100">
          <h2 className="text-xl font-bold text-gray-900">Создать объявление</h2>
          <button 
            onClick={onClose}
            className="p-1 hover:bg-gray-100 rounded transition-colors"
          >
            <X className="w-6 h-6 text-gray-500" />
          </button>
        </div>

        {/* Form Body */}
        <div className="p-6 overflow-y-auto">
          {error && <div className="mb-4 text-red-500 font-bold">{error}</div>}
          <form id="create-announcement-form" onSubmit={handleSubmit} className="space-y-6">
            
            {/* Title */}
            <div>
              <label className="block text-sm font-bold text-gray-700 mb-2">
                Название работы <span className="text-red-500">*</span>
              </label>
              <input 
                type="text" 
                required
                className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none transition-colors"
                placeholder="например, Починить протекающий кран"
                value={formData.title}
                onChange={e => setFormData({...formData, title: e.target.value})}
              />
            </div>

            {/* Category & Budget */}
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-bold text-gray-700 mb-2">Категория</label>
                <div className="relative">
                  <Briefcase className="absolute left-3 top-3.5 w-4 h-4 text-gray-400" />
                  <select 
                    className="w-full border-2 border-gray-300 p-3 pl-10 rounded-md focus:border-gray-800 outline-none appearance-none bg-white"
                    value={formData.categoryId}
                    onChange={e => setFormData({...formData, categoryId: Number(e.target.value)})}
                  >
                    {categories.map(cat => (
                      <option key={cat.id} value={cat.id}>
                        {cat.name}
                      </option>
                    ))}
                  </select>
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-bold text-gray-700 mb-2">Бюджет (примерно)</label>
                <input 
                  type="number" 
                  className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none"
                  placeholder="100"
                  min={1}
                  required
                  value={formData.budget}
                  onChange={e => setFormData({...formData, budget: e.target.value})}
                />
              </div>
            </div>

            {/* Price unit & Schedule */}
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-bold text-gray-700 mb-2">Единица цены</label>
                <select
                  className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none bg-white"
                  value={formData.priceUnitId}
                  onChange={e => setFormData({ ...formData, priceUnitId: Number(e.target.value) })}
                >
                  {priceUnits.map(unit => (
                    <option key={unit.id} value={unit.id}>
                      {unit.name}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-bold text-gray-700 mb-2">Когда актуально</label>
                <input
                  type="text"
                  className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none"
                  placeholder="будни, с 9 до 18"
                  value={formData.schedule}
                  onChange={e => setFormData({ ...formData, schedule: e.target.value })}
                />
              </div>
            </div>

            {/* Location */}
            <div>
              <label className="block text-sm font-bold text-gray-700 mb-2">Местоположение</label>
              <div className="relative">
                <MapPin className="absolute left-3 top-3.5 w-4 h-4 text-gray-400" />
                <input 
                  type="text" 
                  className="w-full border-2 border-gray-300 p-3 pl-10 rounded-md focus:border-gray-800 outline-none"
                  placeholder="ул. Главная, 123, Город"
                  value={formData.location}
                  onChange={e => setFormData({...formData, location: e.target.value})}
                />
              </div>
            </div>

            {/* Description */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <label className="block text-sm font-bold text-gray-700">Описание</label>
                <AIAssistButton 
                  currentText={formData.description}
                  onGenerated={(text) => setFormData({...formData, description: text})}
                  context={`Название работы: ${formData.title}. Локация: ${formData.location}. Бюджет: ${formData.budget}`}
                />
              </div>
              <textarea 
                rows={4}
                className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none resize-none"
                placeholder="Опишите проблему подробно..."
                value={formData.description}
                onChange={e => setFormData({...formData, description: e.target.value})}
              />
            </div>

            {/* Photo Upload Placeholder */}
            <div>
              <label className="block text-sm font-bold text-gray-700 mb-2">Фото</label>
              <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 flex flex-col items-center justify-center text-gray-500 hover:bg-gray-50 hover:border-gray-400 cursor-pointer transition-colors">
                <Upload className="w-8 h-8 mb-2" />
                <span className="text-sm font-medium">Нажмите для загрузки фото</span>
                <span className="text-xs text-gray-400 mt-1">(JPG, PNG макс 5МБ)</span>
              </div>
            </div>

          </form>
        </div>

        {/* Footer */}
        <div className="p-6 border-t-2 border-gray-100 bg-gray-50 flex justify-end gap-3">
          <button 
            onClick={onClose}
            type="button"
            disabled={isLoading}
            className="px-6 py-2.5 font-bold text-gray-600 hover:text-gray-900 hover:bg-gray-200 rounded transition-colors disabled:opacity-50"
          >
            Отмена
          </button>
          <button 
            type="submit"
            form="create-announcement-form"
            disabled={isLoading}
            className="bg-gray-900 text-white px-6 py-2.5 rounded font-bold hover:bg-gray-800 active:translate-y-0.5 transition-all flex items-center gap-2 disabled:opacity-50"
          >
            {isLoading ? 'Загрузка...' : 'Опубликовать объявление'}
          </button>
        </div>
      </div>
    </div>
  );
}
