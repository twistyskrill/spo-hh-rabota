import React, { useState } from 'react';

interface UserProfile {
  name: string;
  email: string;
  phone: string;
  location: string;
  memberSince: string;
}

interface UserProfileEditorProps {
  profile: UserProfile;
  onSave: (profile: UserProfile) => void;
  onBack: () => void;
}

export function UserProfileEditor({ profile, onSave, onBack }: UserProfileEditorProps) {
  const [name, setName] = useState(profile.name);
  const [email, setEmail] = useState(profile.email);
  const [phone, setPhone] = useState(profile.phone);
  const [location, setLocation] = useState(profile.location);

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    onSave({
      ...profile,
      name,
      email,
      phone,
      location,
      memberSince: profile.memberSince || new Date().toLocaleDateString(),
    });
  };

  return (
    <div className="max-w-2xl mx-auto py-8">
      <button onClick={onBack} className="mb-6 flex items-center text-gray-500 hover:text-gray-900 font-medium transition-colors">
        Назад
      </button>
      <div className="bg-white border-2 border-gray-200 rounded-lg p-8 shadow-sm">
        <h2 className="text-2xl font-bold mb-6 text-gray-900">Редактировать профиль пользователя</h2>
        <form onSubmit={handleSave} className="space-y-6">
          <div>
            <label className="block text-sm font-bold text-gray-700 mb-2">Имя</label>
            <input type="text" value={name} onChange={e => setName(e.target.value)} className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none" required />
          </div>
          <div>
            <label className="block text-sm font-bold text-gray-700 mb-2">Email</label>
            <input type="text" value={email} onChange={e => setEmail(e.target.value)} className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none" required />
          </div>
          <div>
            <label className="block text-sm font-bold text-gray-700 mb-2">Телефон</label>
            <input type="text" value={phone} onChange={e => setPhone(e.target.value)} className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none" required />
          </div>
          <div>
            <label className="block text-sm font-bold text-gray-700 mb-2">Локация</label>
            <input type="text" value={location} onChange={e => setLocation(e.target.value)} className="w-full border-2 border-gray-300 p-3 rounded-md focus:border-gray-800 outline-none" required />
          </div>
          <button type="submit" className="w-full bg-gray-800 text-white font-bold py-3 px-4 rounded-md hover:bg-gray-700 transition-all">Сохранить</button>
        </form>
      </div>
    </div>
  );
}
