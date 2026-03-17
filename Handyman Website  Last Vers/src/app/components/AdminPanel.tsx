import React from 'react';
import { ArrowRight } from 'lucide-react';
import { Announcement, Handyman } from '../data';

interface AdminPanelProps {
  pendingAnnouncements: Announcement[];
  pendingHandymen: Handyman[];
  onApproveAnnouncement: (id: number) => void;
  onRejectAnnouncement: (id: number) => void;
  onApproveHandyman: (id: string) => void;
  onRejectHandyman: (id: string) => void;
  onBack: () => void;
}

export function AdminPanel({
  pendingAnnouncements,
  pendingHandymen,
  onApproveAnnouncement,
  onRejectAnnouncement,
  onApproveHandyman,
  onRejectHandyman,
  onBack,
}: AdminPanelProps) {

  return (
    <div className="space-y-8">
      <button
        onClick={onBack}
        className="flex items-center text-gray-500 hover:text-gray-900 font-medium transition-colors"
      >
        <ArrowRight className="w-4 h-4 mr-2 rotate-180" />
        Назад на главную
      </button>

      <div className="border-b border-gray-200 pb-4">
        <h1 className="text-2xl font-bold text-gray-900">Админ-панель</h1>
        <p className="text-gray-500 mt-1">Ручная модерация объявлений и аккаунтов мастеров.</p>
      </div>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-2">Объявления на модерации</h2>
        {pendingAnnouncements.length > 0 ? (
          <ul className="space-y-4">
            {pendingAnnouncements.map((announcement) => (
              <li key={announcement.id} className="bg-white border-2 border-gray-200 rounded-lg p-4">
                <h3 className="font-bold text-lg">{announcement.title}</h3>
                <p className="text-gray-700">{announcement.description || 'Описание не указано'}</p>
                <p className="text-sm text-gray-500 mt-1">Категория: {announcement.category}</p>
                <p className="text-sm text-gray-500">Бюджет: {announcement.budget}</p>
                <div className="mt-2 space-x-2">
                  <button
                    onClick={() => onApproveAnnouncement(announcement.id)}
                    className="bg-gray-800 text-white px-4 py-2 rounded-md hover:bg-gray-700"
                  >
                    Одобрить
                  </button>
                  <button
                    onClick={() => onRejectAnnouncement(announcement.id)}
                    className="border-2 border-gray-800 text-gray-800 px-4 py-2 rounded-md hover:bg-gray-100"
                  >
                    Отклонить
                  </button>
                </div>
              </li>
            ))}
          </ul>
        ) : (
          <p className="text-gray-500">Нет объявлений, ожидающих модерацию.</p>
        )}
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-2">Аккаунты мастеров на модерации</h2>
        {pendingHandymen.length > 0 ? (
          <ul className="space-y-4">
            {pendingHandymen.map((account) => (
              <li key={account.id} className="bg-white border-2 border-gray-200 rounded-lg p-4">
                <h3 className="font-bold text-lg">{account.name || 'Новый мастер'}</h3>
                <p className="text-gray-700">Email: {account.email || 'не указан'}</p>
                <p className="text-gray-700">Специализация: {account.skill}</p>
                <p className="text-gray-700">Ставка: {account.hourlyRate} руб./час</p>
                <div className="mt-2 space-x-2">
                  <button
                    onClick={() => onApproveHandyman(account.id)}
                    className="bg-gray-800 text-white px-4 py-2 rounded-md hover:bg-gray-700"
                  >
                    Одобрить
                  </button>
                  <button
                    onClick={() => onRejectHandyman(account.id)}
                    className="border-2 border-gray-800 text-gray-800 px-4 py-2 rounded-md hover:bg-gray-100"
                  >
                    Отклонить
                  </button>
                </div>
              </li>
            ))}
          </ul>
        ) : (
          <p className="text-gray-500">Нет аккаунтов мастеров, ожидающих модерацию.</p>
        )}
      </section>
    </div>
  );
}