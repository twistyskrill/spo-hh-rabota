export interface Review {
  id: string;
  author: string;
  rating: number;
  text: string;
  date: string;
}

export interface Handyman {
  id: string;
  name: string;
  skill: string;
  rating: number;
  hourlyRate: number;
  description: string;
  reviews: Review[];
  avatarUrl?: string;
  email?: string; // Added optional email property to match usage in App.tsx
}

export const MOCK_HANDYMEN: Handyman[] = [
  {
    id: "1",
    name: "Иван Петров",
    skill: "Сантехника",
    rating: 4.8,
    hourlyRate: 50,
    description: "Опытный сантехник с 10-летним стажем. Я занимаюсь всем: от устранения протечек до полной установки.",
    reviews: [
      { id: "r1", author: "Анна М.", rating: 5, text: "Быстро и надежно!", date: "2023-10-15" },
      { id: "r2", author: "Иван Д.", rating: 4, text: "Хорошая работа, немного опоздал.", date: "2023-09-20" }
    ]
  },
  {
    id: "2",
    name: "Сергей Кузнецов",
    skill: "Столярные работы",
    rating: 4.9,
    hourlyRate: 65,
    description: "Пользовательская мебель, ремонт шкафов и структурные деревянные работы.",
    reviews: [
      { id: "r3", author: "Михаил Т.", rating: 5, text: "Потрясающее мастерство.", date: "2023-11-01" }
    ]
  },
  {
    id: "3",
    name: "Елена Смирнова",
    skill: "Электрика",
    rating: 4.7,
    hourlyRate: 70,
    description: "Лицензированный электрик для ремонта и модернизации жилых помещений.",
    reviews: [
      { id: "r4", author: "Лидия К.", rating: 5, text: "Исправил мою проводку за час.", date: "2023-10-05" }
    ]
  },
  {
    id: "4",
    name: "Дмитрий Васильев",
    skill: "Ландшафтный дизайн",
    rating: 4.5,
    hourlyRate: 40,
    description: "Косилка, обрезка и общий уход за садом.",
    reviews: [
      { id: "r5", author: "Тимофей Х.", rating: 4, text: "Солидная работа.", date: "2023-09-12" }
    ]
  },
  {
    id: "5",
    name: "Ольга Никитина",
    skill: "Малярные работы",
    rating: 4.6,
    hourlyRate: 45,
    description: "Услуги внутренней и внешней покраски. Чисто и точно.",
    reviews: []
  }
];

export interface Announcement {
  id: number;
  title: string;
  category: string;
  status: string;
  date: string;
  budget: string;
  handyman: string | null;
  location?: string;
  description?: string;
}

export const MOCK_ANNOUNCEMENTS: Announcement[] = [
  {
    id: 1,
    title: 'Починить протечку в кухонной раковине',
    category: 'Сантехника',
    status: 'В работе',
    date: '2023-10-25',
    budget: '150 руб.',
    handyman: 'Иван Петров',
    location: 'ул. Главная, 123',
    description: 'Раковина протекает из нижней трубы.'
  },
  {
    id: 2,
    title: 'Установить потолочный вентилятор',
    category: 'Электрика',
    status: 'Открыто',
    date: '2023-10-28',
    budget: '80 руб.',
    handyman: null,
    location: 'ул. Дубовая, 456',
    description: 'Нужно установить новый вентилятор в спальне.'
  }
];
