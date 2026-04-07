export const VALIDATION = {
  // Allows Cyrillic/Latin letters, spaces, hyphen, apostrophe. 2..60 chars.
  name: /^[A-Za-zА-Яа-яЁё\s\-']{2,60}$/,

  // Basic email (HTML type="email" also validates, but pattern keeps rules consistent).
  email: /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/,

  // Phone: +7 (999) 123-45-67, 89991234567, etc. 7..20 chars.
  phone: /^\+?[0-9\s()\-]{7,20}$/,

  // Location/address: letters/numbers/spaces and common punctuation. 2..120 chars.
  location: /^[A-Za-zА-Яа-яЁё0-9\s,\.\-]{2,120}$/,

  // Password: at least 6 chars, allow common printable symbols without spaces.
  password: /^[^\s]{6,64}$/,

  // Announcement title: letters/numbers/spaces and basic punctuation. 3..80 chars.
  title: /^[A-Za-zА-Яа-яЁё0-9\s,\.\-"'()!?:]{3,80}$/,

  // Schedule/free text short field: 0..80 chars, no angle brackets.
  schedule: /^[^<>]{0,80}$/,
};
