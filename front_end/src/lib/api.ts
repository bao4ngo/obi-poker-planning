const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export async function createSession(name: string, hostName: string): Promise<{ sessionId: string; hostId: string }> {
  const response = await fetch(`${API_BASE_URL}/api/sessions`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name, hostName }),
  });

  if (!response.ok) {
    throw new Error('Failed to create session');
  }

  return response.json();
}

export async function getSession(sessionId: string) {
  const response = await fetch(`${API_BASE_URL}/api/sessions/${sessionId}`);

  if (!response.ok) {
    throw new Error('Failed to get session');
  }

  return response.json();
}

export async function addItem(sessionId: string, title: string, description: string) {
  const response = await fetch(`${API_BASE_URL}/api/sessions/${sessionId}/items`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ title, description }),
  });

  if (!response.ok) {
    throw new Error('Failed to add item');
  }

  return response.json();
}

export async function setCurrentItem(sessionId: string, itemId: string) {
  const response = await fetch(`${API_BASE_URL}/api/sessions/${sessionId}/current-item`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ itemId }),
  });

  if (!response.ok) {
    throw new Error('Failed to set current item');
  }

  return response.json();
}

export function connectWebSocket(sessionId: string): WebSocket {
  const wsUrl = API_BASE_URL.replace(/^http/, 'ws').replace(/^https/, 'wss');
  return new WebSocket(`${wsUrl}/ws/${sessionId}`);
}
