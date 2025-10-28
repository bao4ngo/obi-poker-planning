export interface User {
  id: string;
  name: string;
  isHost: boolean;
  vote?: string;
  connected: boolean;
}

export interface PlanningItem {
  id: string;
  title: string;
  description: string;
  votes: { [userId: string]: string };
  revealed: boolean;
  finalEstimate?: string;
}

export interface Session {
  id: string;
  name: string;
  hostId: string;
  users: { [userId: string]: User };
  items: PlanningItem[];
  currentItemId?: string;
  createdAt: string;
}

export interface WSMessage {
  type: string;
  payload: any;
}

export const CARD_VALUES = ['0', '1', '2', '3', '5', '8', '13', '21', '34', '55', '89', '?'];
