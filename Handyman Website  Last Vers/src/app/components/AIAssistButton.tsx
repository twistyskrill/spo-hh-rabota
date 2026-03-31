import React, { useState } from 'react';
import { Sparkles, Loader2 } from 'lucide-react';

interface AIAssistButtonProps {
  currentText: string;
  onGenerated: (text: string) => void;
  context?: string;
}

export function AIAssistButton({ currentText, onGenerated, context = '' }: AIAssistButtonProps) {
  const [isLoading, setIsLoading] = useState(false);

  const handleGenerate = async () => {
    setIsLoading(true);
    try {
      const prompt = currentText.trim().length > 0
        ? `Улучши следующее описание (сделай его более профессиональным, понятным и привлекательным, исправь ошибки). Верни только улучшенный текст на русском языке, без вводных слов и комментариев. Дополнительный контекст: ${context}\n\nОписание:\n${currentText}`
        : `Придумай хорошее, подробное и привлекательное описание. Верни только текст описания на русском языке, без дополнительных комментариев и вводных слов. Дополнительный контекст: ${context}`;

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
