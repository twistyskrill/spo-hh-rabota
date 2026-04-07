import React, { useState } from 'react';
import { Sparkles, Loader2 } from 'lucide-react';

interface AIAssistButtonProps {
  currentText: string;
  onGenerated: (text: string) => void;
  context?: string;
  userType?: 'worker' | 'user';
}

export function AIAssistButton({ currentText, onGenerated, context = '', userType = 'worker' }: AIAssistButtonProps) {
  const [isLoading, setIsLoading] = useState(false);

  const handleGenerate = async () => {
    setIsLoading(true);
    try {
      let prompt = '';
      if (userType === 'worker') {
        prompt = currentText.trim().length > 0
          ? `Улучши следующее описание профиля мастера (сделай его более профессиональным, понятным и привлекательным для клиентов, исправь ошибки). Верни только улучшенный текст на русском языке, без вводных слов и комментариев. Текст должен быть на русском языке. Дополнительный контекст: ${context}\n\nОписание:\n${currentText}`
          : `Создай профессиональное, официальное описание для профиля мастера по ремонту. Текст должен быть грамотным, понятным и привлекательным для потенциальных клиентов. Укажи, что мастер квалифицированный, надежный, ответственный и работает по честным ценам. Исключи разговорные выражения. Длина текста — примерно 3–5 предложений. Верни только текст описания. Текст должен быть на русском языке. Дополнительный контекст: ${context}`;
      } else {
        prompt = currentText.trim().length > 0
          ? `Улучши следующее описание задачи/заказа от лица заказчика (сделай его более понятным, детальным и привлекательным для исполнителей, исправь ошибки). Верни только улучшенный текст на русском языке, без вводных слов и комментариев. Текст должен быть на русском языке. Дополнительный контекст: ${context}\n\nОписание:\n${currentText}`
          : `Создай четкое и понятное описание задачи/заказа для мастера по ремонту от лица клиента. Описание должно привлекать хороших специалистов, быть информативным и описывать суть проблемы или задачи. Длина текста — примерно 3–5 предложений. Верни только текст описания на русском языке, без вводных слов. Текст должен быть на русском языке. Дополнительный контекст: ${context}`;
      }

      const res = await fetch('http://localhost:11434/api/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          model: 'llama3:latest',
          prompt: prompt,
          stream: false,
        }),
      });

      if (!res.ok) {
        throw new Error('Оллома вернула ошибку');
      }

      const data = await res.json();
      if (data.response) {
        onGenerated(data.response.trim());
      }
    } catch (err) {
      console.error(err);
      alert('Ошибка при обращении к Ollama. Убедитесь, что сервер запущен (ollama serve).');
    } finally {
      setIsLoading(false);
    }
  };

  const isImprove = currentText.trim().length > 0;

  return (
    <button
      type="button"
      onClick={handleGenerate}
      disabled={isLoading}
      className={`flex items-center gap-2 px-3 py-1.5 text-sm font-bold text-white rounded-md transition-colors disabled:opacity-50 ${
        isImprove ? 'bg-indigo-600 hover:bg-indigo-700' : 'bg-pink-600 hover:bg-pink-700'
      }`}
      title="Сгенерировать с помощью AI (Ollama)"
    >
      {isLoading ? (
        <Loader2 className="w-4 h-4 animate-spin" />
      ) : (
        <Sparkles className="w-4 h-4" />
      )}
      {isLoading ? 'Думает...' : isImprove ? 'Улучшить' : 'Придумать'}
    </button>
  );
}
