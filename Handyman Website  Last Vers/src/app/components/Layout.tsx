import React, { useState } from 'react';
import { Search, User, Sun, Moon, LogOut } from 'lucide-react';
import { VALIDATION } from '../validation';

interface LayoutProps {
  children: React.ReactNode;
  onNavigate: (page: string) => void;
  currentPage: string;
  onLogout?: () => void;
  isAdmin?: boolean;
  onOpenAdminPanel?: () => void;
}

export function Layout({ children, onNavigate, currentPage, onLogout, isAdmin = false, onOpenAdminPanel }: LayoutProps) {
  const [isDarkMode, setIsDarkMode] = useState(false);

  const toggleTheme = () => {
    setIsDarkMode(!isDarkMode);
    document.documentElement.classList.toggle('dark');
  };

  return (
    <div className={`min-h-screen flex flex-col font-sans text-gray-900 bg-gray-50 ${isDarkMode ? 'dark:bg-gray-900 dark:text-gray-100' : ''}`}>
      {/* Header */}
      <header className="border-b-2 border-gray-800 bg-white p-4 sticky top-0 z-10 flex items-center justify-between shadow-sm">
        <div 
          className="flex items-center gap-2 cursor-pointer" 
          onClick={() => onNavigate('home')}
        >
          <div className="w-8 h-8 bg-gray-800 rounded flex items-center justify-center text-white font-bold">H</div>
          <span className="font-bold text-xl tracking-tight">HandyApp</span>
        </div>

        <div className="flex-1 max-w-md mx-4 hidden md:flex items-center border-2 border-gray-300 rounded-md px-3 py-1.5 bg-gray-50 focus-within:border-gray-800 focus-within:bg-white transition-colors">
          <Search className="w-5 h-5 text-gray-500 mr-2" />
          <input 
            type="text" 
            placeholder="Поиск мастера..." 
            className="bg-transparent border-none outline-none w-full text-sm"
          />
        </div>

        <div className="flex items-center gap-4">
          {isAdmin && onOpenAdminPanel && (
            <button
              onClick={onOpenAdminPanel}
              className={`hidden md:inline-flex items-center gap-2 px-3 py-1.5 border-2 rounded-full text-sm font-bold transition-colors ${
                currentPage === 'admin'
                  ? 'bg-gray-800 text-white border-gray-800'
                  : 'border-gray-800 text-gray-800 hover:bg-gray-100'
              }`}
            >
              Админ-панель
            </button>
          )}

          <button 
            onClick={toggleTheme}
            className="p-2 border-2 border-transparent hover:border-gray-200 rounded-full transition-all"
            aria-label="Toggle Theme"
          >
            {isDarkMode ? <Sun className="w-5 h-5" /> : <Moon className="w-5 h-5" />}
          </button>

          {onLogout && (
            <button 
              onClick={onLogout}
              className="hidden md:inline-flex items-center gap-2 px-3 py-1.5 border-2 border-gray-800 rounded-full text-sm font-bold text-gray-800 hover:bg-gray-100 transition-colors"
            >
              <LogOut className="w-4 h-4" />
              Выйти
            </button>
          )}
          
          <button 
            onClick={() => onNavigate('profile-user')}
            className={`p-2 border-2 rounded-full transition-colors ${
              currentPage === 'profile-user' 
                ? 'bg-gray-800 text-white border-gray-800' 
                : 'border-gray-800 text-gray-800 hover:bg-gray-100'
            }`}
            aria-label="User Profile"
          >
            <User className="w-5 h-5" />
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 container mx-auto p-4 md:p-8 max-w-5xl">
        {children}
      </main>

      {/* Footer */}
      <footer className="border-t-2 border-gray-800 bg-white p-8 mt-auto">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
          <div>
            <h4 className="font-bold mb-4 uppercase text-sm tracking-wider">Компания</h4>
            <ul className="space-y-2 text-sm text-gray-600">
              <li className="hover:underline cursor-pointer">О нас</li>
              <li className="hover:underline cursor-pointer">Карьера</li>
              <li className="hover:underline cursor-pointer">Пресса</li>
            </ul>
          </div>
          <div>
            <h4 className="font-bold mb-4 uppercase text-sm tracking-wider">Поддержка</h4>
            <ul className="space-y-2 text-sm text-gray-600">
              <li className="hover:underline cursor-pointer">Контакты</li>
              <li className="hover:underline cursor-pointer">Центр помощи</li>
              <li className="hover:underline cursor-pointer">Безопасность</li>
            </ul>
          </div>
          <div>
            <h4 className="font-bold mb-4 uppercase text-sm tracking-wider">Юридическая информация</h4>
            <ul className="space-y-2 text-sm text-gray-600">
              <li className="hover:underline cursor-pointer">Условия обслуживания</li>
              <li className="hover:underline cursor-pointer">Политика конфиденциальности</li>
              <li className="hover:underline cursor-pointer">Политика cookies</li>
            </ul>
          </div>
          <div>
            <h4 className="font-bold mb-4 uppercase text-sm tracking-wider">Подписка</h4>
            <div className="flex gap-2">
              <input type="email" placeholder="Электронная почта" className="border-2 border-gray-300 p-2 w-full text-sm rounded-sm" pattern={VALIDATION.email.source} />
              <button className="bg-gray-800 text-white px-4 py-2 font-bold text-sm uppercase rounded-sm">Отправить</button>
            </div>
          </div>
        </div>
        <div className="mt-8 pt-8 border-t border-gray-200 text-center text-xs text-gray-500">
          © {new Date().getFullYear()} HandyApp Inc. Все права защищены.
        </div>
      </footer>
    </div>
  );
}
