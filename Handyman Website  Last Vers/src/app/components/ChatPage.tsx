import React, { useState } from 'react';
import { Send, ArrowLeft, Phone, Video } from 'lucide-react';
import { MOCK_HANDYMEN } from '../data';

interface ChatPageProps {
  handymanId: string;
  onBack: () => void;
}

export function ChatPage({ handymanId, onBack }: ChatPageProps) {
  const handyman = MOCK_HANDYMEN.find(h => h.id === handymanId);
  const [message, setMessage] = useState('');
  const [messages, setMessages] = useState([
    { id: 1, text: "Привет, меня интересуют ваши услуги.", sender: 'user', time: '10:00' },
    { id: 2, text: `Привет! Я рад помочь. Что нужно сделать?`, sender: 'handyman', time: '10:05' }
  ]);

  if (!handyman) return <div>Мастер не найден</div>;

  const handleSend = (e: React.FormEvent) => {
    e.preventDefault();
    if (!message.trim()) return;
    
    setMessages([...messages, { 
      id: Date.now(), 
      text: message, 
      sender: 'user', 
      time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) 
    }]);
    setMessage('');

    // Simulate reply
    setTimeout(() => {
        setMessages(prev => [...prev, {
            id: Date.now() + 1,
            text: "Спасибо за детали. Я могу зайти завтра.",
            sender: 'handyman',
            time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
        }]);
    }, 1500);
  };

  return (
    <div className="h-[calc(100vh-140px)] flex flex-col bg-white border-2 border-gray-200 rounded-lg shadow-sm overflow-hidden">
      {/* Chat Header */}
      <div className="p-4 border-b-2 border-gray-100 flex items-center justify-between bg-gray-50">
        <div className="flex items-center gap-3">
          <button onClick={onBack} className="p-1 hover:bg-gray-200 rounded-full transition-colors md:hidden">
            <ArrowLeft className="w-5 h-5" />
          </button>
          <div className="w-10 h-10 bg-gray-300 rounded-full flex items-center justify-center font-bold text-gray-600">
             {handyman.name.charAt(0)}
          </div>
          <div>
            <h3 className="font-bold text-gray-900">{handyman.name}</h3>
            <span className="text-xs text-green-600 font-medium flex items-center gap-1">
              <span className="w-1.5 h-1.5 bg-green-500 rounded-full"></span>
              Онлайн
            </span>
          </div>
        </div>
        <div className="flex gap-2">
            <button className="p-2 text-gray-500 hover:bg-gray-200 rounded-full">
                <Phone className="w-5 h-5" />
            </button>
            <button className="p-2 text-gray-500 hover:bg-gray-200 rounded-full">
                <Video className="w-5 h-5" />
            </button>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4 bg-gray-50/50">
        {messages.map(msg => (
          <div 
            key={msg.id} 
            className={`flex flex-col ${msg.sender === 'user' ? 'items-end' : 'items-start'}`}
          >
            <div 
              className={`max-w-[80%] p-3 rounded-lg text-sm ${
                msg.sender === 'user' 
                  ? 'bg-gray-800 text-white rounded-br-none' 
                  : 'bg-white border border-gray-200 text-gray-800 rounded-bl-none shadow-sm'
              }`}
            >
              {msg.text}
            </div>
            <span className="text-[10px] text-gray-400 mt-1 px-1">
              {msg.time}
            </span>
          </div>
        ))}
      </div>

      {/* Input */}
      <form onSubmit={handleSend} className="p-4 border-t-2 border-gray-100 bg-white flex gap-2">
        <input 
          type="text" 
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Type your message..." 
          className="flex-1 border-2 border-gray-200 rounded-md px-4 py-2 focus:border-gray-800 outline-none transition-colors"
          minLength={1}
          maxLength={500}
        />
        <button 
          type="submit" 
          disabled={!message.trim()}
          className="bg-gray-800 text-white p-2 rounded-md hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <Send className="w-5 h-5" />
        </button>
      </form>
    </div>
  );
}
