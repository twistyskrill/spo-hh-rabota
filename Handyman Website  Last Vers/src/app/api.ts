export const API_URL = 'http://localhost:8080';

const getErrorMessage = async (res: Response, fallback: string): Promise<string> => {
  const text = await res.text();
  if (!text) return fallback;

  try {
    const parsed = JSON.parse(text);
    if (parsed?.message && typeof parsed.message === 'string') return parsed.message;
  } catch {
    // Ignore JSON parse errors and return raw text below.
  }

  return text;
};

export const getAuthHeaders = (): Record<string, string> => {
  const token = localStorage.getItem('token');
  return token ? { Authorization: `Bearer ${token}` } : {};
};

export const api = {
  // Auth
  register: async (data: any) => {
    const res = await fetch(`${API_URL}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error(await getErrorMessage(res, 'Registration failed'));
    return res.json();
  },
  login: async (data: any) => {
    const res = await fetch(`${API_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error(await getErrorMessage(res, 'Login failed'));
    return res.json();
  },

  // Profile
  getProfile: async () => {
    const res = await fetch(`${API_URL}/profile`, {
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to get profile');
    return res.json();
  },
  updateProfile: async (data: any) => {
    const res = await fetch(`${API_URL}/profile`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to update profile');
    return res.json();
  },

  // Ads
  getAds: async (limit = 10, offset = 0) => {
    const res = await fetch(`${API_URL}/ads?limit=${limit}&offset=${offset}`);
    if (!res.ok) throw new Error('Failed to fetch ads');
    return res.json();
  },
  getMyAds: async () => {
    const res = await fetch(`${API_URL}/my-ads`, {
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch my ads');
    return res.json();
  },
  getAdById: async (id: number) => {
    const res = await fetch(`${API_URL}/my-ads/${id}`, {
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch ad');
    return res.json();
  },
  getCategories: async () => {
    const res = await fetch(`${API_URL}/info/categories`);
    if (!res.ok) throw new Error('Failed to fetch categories');
    return res.json();
  },
  getPriceUnits: async () => {
    const res = await fetch(`${API_URL}/info/price_units`);
    if (!res.ok) throw new Error('Failed to fetch price units');
    return res.json();
  },
  createAd: async (data: any) => {
    const res = await fetch(`${API_URL}/my-ads/`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to create ad');
    return res.json();
  },
  updateAd: async (id: number, data: any) => {
    const res = await fetch(`${API_URL}/my-ads/${id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to update ad');
    return res.json();
  },
  deleteAd: async (id: number) => {
    const res = await fetch(`${API_URL}/my-ads/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to delete ad');
    return res.json();
  },

  // Admin - ads moderation
  getAdminAds: async (status: string = 'pending') => {
    const query = status ? `?status=${encodeURIComponent(status)}` : '';
    const res = await fetch(`${API_URL}/admin/ads${query}`, {
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch admin ads');
    return res.json();
  },
  approveAd: async (id: number) => {
    const res = await fetch(`${API_URL}/admin/ads/${id}/approve`, {
      method: 'PATCH',
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to approve ad');
    return res.json();
  },
  rejectAd: async (id: number) => {
    const res = await fetch(`${API_URL}/admin/ads/${id}/reject`, {
      method: 'PATCH',
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to reject ad');
    return res.json();
  },

  // Handymen
  getHandymen: async (limit = 10, offset = 0) => {
    const res = await fetch(`${API_URL}/handyman?limit=${limit}&offset=${offset}`);
    if (!res.ok) throw new Error('Failed to fetch handymen');
    return res.json();
  },
  getHandymanById: async (id: number) => {
    const res = await fetch(`${API_URL}/handyman/${id}`);
    if (!res.ok) throw new Error('Failed to fetch handyman');
    return res.json();
  },
  getHandymanCategories: async () => {
    const res = await fetch(`${API_URL}/handyman/categories`, {
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch handyman categories');
    return res.json();
  },
  addHandymanCategories: async (data: { category_names: string[] }) => {
    const res = await fetch(`${API_URL}/handyman/categories`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to add handyman categories');
    return res.json();
  },
  removeHandymanCategories: async (data: { category_names: string[] }) => {
    const res = await fetch(`${API_URL}/handyman/categories`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to remove handyman categories');
    return res.json();
  },

  // Responses (Отклики)
  getResponses: async () => {
    const res = await fetch(`${API_URL}/responses`, {
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch responses');
    return res.json();
  },
  createResponse: async (data: any) => {
    const res = await fetch(`${API_URL}/responses`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to create response');
    return res.json();
  },
  deleteResponse: async (id: number) => {
    const res = await fetch(`${API_URL}/responses/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to delete response');
    return res.json();
  },

  // Admin - worker moderation
  getAdminWorkers: async (status: string = 'pending') => {
    const query = status ? `?status=${encodeURIComponent(status)}` : '';
    const res = await fetch(`${API_URL}/admin/workers${query}`, {
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to fetch admin workers');
    return res.json();
  },
  approveWorker: async (workerId: number) => {
    const res = await fetch(`${API_URL}/admin/workers/${workerId}/approve`, {
      method: 'PATCH',
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to approve worker');
    return res.json();
  },
  rejectWorker: async (workerId: number) => {
    const res = await fetch(`${API_URL}/admin/workers/${workerId}/reject`, {
      method: 'PATCH',
      headers: getAuthHeaders(),
    });
    if (!res.ok) throw new Error('Failed to reject worker');
    return res.json();
  },
};
