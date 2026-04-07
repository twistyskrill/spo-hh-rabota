import React, { useState, useEffect } from 'react';
import { Handyman } from '../data';
import { api } from '../api';
import { AIAssistButton } from './AIAssistButton';
import { VALIDATION } from '../validation';

interface HandymanProfileEditorProps {
  handyman: Handyman;
  onSave: (handyman: Handyman) => void;
  onBack: () => void;
}

export function HandymanProfileEditor({ handyman, onSave, onBack }: HandymanProfileEditorProps) {
  const [name, setName] = useState(handyman.name);
  const [skill, setSkill] = useState(handyman.skill || '');
  const [hourlyRate, setHourlyRate] = useState(handyman.hourlyRate);
  const [description, setDescription] = useState(handyman.description);

  // Categories from DB
  const [categories, setCategories] = useState<{ id: number; name: string }[]>([]);

  useEffect(() => {
    api.getCategories().then((data) => {
      if (Array.isArray(data)) {
        setCategories(data);
        // If current skill is empty or not in the list, set to first available
        if (data.length > 0 && !skill) {
          setSkill(data[0].name);
        }
      }
    }).catch(console.error);
  }, []);

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    onSave({ ...handyman, name, skill, hourlyRate, description });
  };

  return (
    <div className="max-w-2xl mx-auto py-8">
      <button onClick={onBack} className="mb-6 flex items-center text-gray-500 hover:text-gray-900 font-medium transition-colors">
        Назад
      </button>
      <div className="bg-white border-2 border-gray-200 rounded-lg p-8 shadow-sm">
        <h2 className="text-2xl font-bold mb-6 text-gray-900">Редактировать профиль мастера</h2>
        <form onSubmit={handleSave} className="space-y-6">
          <div>
            <label className="block text-sm font-bold text-gray-700 mb-2">Имя</label>
            <input type="text" value={name} onChange={e => setName(e.target.value)} className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none" pattern={VALIDATION.name.source} required />
          </div>
          <div>
            <label className="block text-sm font-bold text-gray-700 mb-2">Специализация</label>
            <select value={skill} onChange={e => setSkill(e.target.value)} className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none">
              {categories.length > 0 ? (
                categories.map((cat) => (
                  <option key={cat.id} value={cat.name}>{cat.name}</option>
                ))
              ) : (
                <option disabled>Загрузка...</option>
              )}
            </select>
          </div>
          <div>
            <label className="block text-sm font-bold text-gray-700 mb-2">Почасовая ставка (руб.)</label>
            <input type="number" value={hourlyRate} onChange={e => setHourlyRate(Number(e.target.value))} className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none" min={1} max={1000000} required />
          </div>
          <div>
            <div className="flex items-center justify-between mb-2">
              <label className="block text-sm font-bold text-gray-700">Описание</label>
              <AIAssistButton 
                currentText={description}
                onGenerated={setDescription}
                context={`Имя мастера: ${name}. Специализация: ${skill}.`}
              />
            </div>
            <textarea value={description} onChange={e => setDescription(e.target.value)} className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none resize-none" rows={3} minLength={10} maxLength={2000} required />
          </div>
          <button type="submit" className="w-full bg-gray-800 text-white font-bold py-3 px-4 rounded-md hover:bg-gray-700 transition-all">Сохранить</button>
        </form>
      </div>
    </div>
  );
}
