import { useEffect, useState, useCallback } from 'react';
import { useRouter } from 'next/router';
import { Session, PlanningItem, User, WSMessage, CARD_VALUES } from '@/types';
import { connectWebSocket, addItem, setCurrentItem } from '@/lib/api';

export default function SessionPage() {
  const router = useRouter();
  const { sessionId, userId, userName } = router.query;

  const [session, setSession] = useState<Session | null>(null);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const [selectedVote, setSelectedVote] = useState<string | null>(null);
  const [connectionError, setConnectionError] = useState<string | null>(null);

  // Form states
  const [newItemTitle, setNewItemTitle] = useState('');
  const [newItemDescription, setNewItemDescription] = useState('');
  const [showAddItem, setShowAddItem] = useState(false);

  useEffect(() => {
    if (!sessionId || !userName) return;

    const websocket = connectWebSocket(sessionId as string);

    websocket.onopen = () => {
      console.log('WebSocket connected');
      websocket.send(JSON.stringify({
        userName: userName,
        userId: userId || '',
      }));
    };

    websocket.onmessage = (event) => {
      const message: WSMessage = JSON.parse(event.data);
      handleWebSocketMessage(message);
    };

    websocket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    websocket.onclose = () => {
      console.log('WebSocket disconnected');
      setConnected(false);
    };

    setWs(websocket);

    return () => {
      websocket.close();
    };
  }, [sessionId, userName, userId]);

  const handleWebSocketMessage = useCallback((message: WSMessage) => {
    console.log('Received message:', message);

    switch (message.type) {
      case 'error':
        setConnectionError(message.payload.error);
        setConnected(false);
        // Redirect back to home page with error and keep the join link
        setTimeout(() => {
          router.push(`/?join=${sessionId}&error=${encodeURIComponent(message.payload.error)}`);
        }, 2000);
        break;

      case 'welcome':
        setConnected(true);
        setConnectionError(null);
        setSession(message.payload.session);
        const user = message.payload.session.users[message.payload.userId];
        setCurrentUser(user);
        break;

      case 'user_joined':
        setSession((prev) => {
          if (!prev) return prev;
          return {
            ...prev,
            users: {
              ...prev.users,
              [message.payload.id]: message.payload,
            },
          };
        });
        break;

      case 'user_left':
        setSession((prev) => {
          if (!prev) return prev;
          const newUsers = { ...prev.users };
          delete newUsers[message.payload.userId];
          return {
            ...prev,
            users: newUsers,
          };
        });
        break;

      case 'item_added':
        setSession((prev) => {
          if (!prev) return prev;
          return {
            ...prev,
            items: [...prev.items, message.payload],
          };
        });
        break;

      case 'current_item_changed':
        setSession((prev) => {
          if (!prev) return prev;
          return {
            ...prev,
            currentItemId: message.payload.itemId,
          };
        });
        setSelectedVote(null);
        break;

      case 'vote_submitted':
        setSession((prev) => {
          if (!prev) return prev;
          const items = prev.items.map((item) => {
            if (item.id === message.payload.itemId) {
              return {
                ...item,
                votes: {
                  ...item.votes,
                  [message.payload.userId]: message.payload.hasVoted ? '✓' : '',
                },
              };
            }
            return item;
          });
          return { ...prev, items };
        });
        break;

      case 'votes_revealed':
        setSession((prev) => {
          if (!prev) return prev;
          const items = prev.items.map((item) => {
            if (item.id === message.payload.id) {
              return message.payload;
            }
            return item;
          });
          return { ...prev, items };
        });
        break;

      case 'votes_reset':
        setSession((prev) => {
          if (!prev) return prev;
          const items = prev.items.map((item) => {
            if (item.id === message.payload.itemId) {
              return {
                ...item,
                votes: {},
                revealed: false,
              };
            }
            return item;
          });
          return { ...prev, items };
        });
        setSelectedVote(null);
        break;

      case 'final_estimate_set':
        setSession((prev) => {
          if (!prev) return prev;
          const items = prev.items.map((item) => {
            if (item.id === message.payload.itemId) {
              return {
                ...item,
                finalEstimate: message.payload.estimate,
              };
            }
            return item;
          });
          return { ...prev, items };
        });
        break;
    }
  }, []);

  const handleVote = (vote: string) => {
    if (!ws || !session?.currentItemId) return;

    setSelectedVote(vote);
    ws.send(JSON.stringify({
      type: 'vote',
      payload: {
        itemId: session.currentItemId,
        vote: vote,
      },
    }));
  };

  const handleRevealVotes = () => {
    if (!ws || !session?.currentItemId) return;

    ws.send(JSON.stringify({
      type: 'reveal_votes',
      payload: {
        itemId: session.currentItemId,
      },
    }));
  };

  const handleResetVotes = () => {
    if (!ws || !session?.currentItemId) return;

    ws.send(JSON.stringify({
      type: 'reset_votes',
      payload: {
        itemId: session.currentItemId,
      },
    }));
  };

  const handleSetFinalEstimate = (estimate: string) => {
    if (!ws || !session?.currentItemId) return;

    ws.send(JSON.stringify({
      type: 'set_final_estimate',
      payload: {
        itemId: session.currentItemId,
        estimate: estimate,
      },
    }));
  };

  const handleAddItem = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!sessionId || !newItemTitle.trim()) return;

    try {
      await addItem(sessionId as string, newItemTitle, newItemDescription);
      setNewItemTitle('');
      setNewItemDescription('');
      setShowAddItem(false);
    } catch (error) {
      console.error('Failed to add item:', error);
    }
  };

  const handleSelectItem = async (itemId: string) => {
    if (!sessionId) return;

    try {
      await setCurrentItem(sessionId as string, itemId);
    } catch (error) {
      console.error('Failed to set current item:', error);
    }
  };

  const getCurrentItem = (): PlanningItem | null => {
    if (!session?.currentItemId) return null;
    return session.items.find((item) => item.id === session.currentItemId) || null;
  };

  const copySessionLink = () => {
    const link = `${window.location.origin}/?join=${sessionId}`;
    navigator.clipboard.writeText(link);
    alert('Session link copied to clipboard!');
  };

  if (!connected || !session || !currentUser) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          {connectionError ? (
            <div className="bg-red-50 border-2 border-red-200 rounded-lg p-6 max-w-md">
              <div className="text-red-600 text-5xl mb-4">❌</div>
              <h2 className="text-xl font-bold text-red-900 mb-2">Connection Error</h2>
              <p className="text-red-700 mb-4">{connectionError}</p>
              <p className="text-sm text-gray-600">Redirecting you back...</p>
            </div>
          ) : (
            <>
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
              <p className="text-gray-600">Connecting to session...</p>
            </>
          )}
        </div>
      </div>
    );
  }

  const currentItem = getCurrentItem();
  const isHost = currentUser.isHost;

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">{session.name}</h1>
              <p className="text-sm text-gray-600">
                Logged in as: <span className="font-semibold">{currentUser.name}</span>
                {isHost && <span className="ml-2 text-blue-600">(Host)</span>}
              </p>
            </div>
            <button
              onClick={copySessionLink}
              className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition text-sm"
            >
              📋 Share Session
            </button>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-8 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Column: Items List */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-xl font-semibold">Planning Items</h2>
                {isHost && (
                  <button
                    onClick={() => setShowAddItem(!showAddItem)}
                    className="bg-green-600 text-white px-3 py-1 rounded text-sm hover:bg-green-700 transition"
                  >
                    + Add
                  </button>
                )}
              </div>

              {showAddItem && isHost && (
                <form onSubmit={handleAddItem} className="mb-4 p-4 bg-gray-50 rounded-lg">
                  <input
                    type="text"
                    value={newItemTitle}
                    onChange={(e) => setNewItemTitle(e.target.value)}
                    placeholder="Item title"
                    className="w-full px-3 py-2 border border-gray-300 rounded mb-2 focus:ring-2 focus:ring-blue-500 outline-none"
                    required
                  />
                  <textarea
                    value={newItemDescription}
                    onChange={(e) => setNewItemDescription(e.target.value)}
                    placeholder="Description (optional)"
                    className="w-full px-3 py-2 border border-gray-300 rounded mb-2 focus:ring-2 focus:ring-blue-500 outline-none"
                    rows={2}
                  />
                  <div className="flex gap-2">
                    <button
                      type="submit"
                      className="flex-1 bg-blue-600 text-white px-3 py-2 rounded hover:bg-blue-700 transition text-sm"
                    >
                      Add Item
                    </button>
                    <button
                      type="button"
                      onClick={() => setShowAddItem(false)}
                      className="px-3 py-2 border border-gray-300 rounded hover:bg-gray-50 transition text-sm"
                    >
                      Cancel
                    </button>
                  </div>
                </form>
              )}

              <div className="space-y-2">
                {session.items.map((item) => (
                  <div
                    key={item.id}
                    onClick={() => isHost && handleSelectItem(item.id)}
                    className={`p-3 rounded-lg border-2 transition cursor-pointer ${
                      session.currentItemId === item.id
                        ? 'border-blue-500 bg-blue-50'
                        : 'border-gray-200 hover:border-gray-300'
                    }`}
                  >
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <h3 className="font-semibold text-gray-900">{item.title}</h3>
                        {item.description && (
                          <p className="text-sm text-gray-600 mt-1">{item.description}</p>
                        )}
                      </div>
                      {item.finalEstimate && (
                        <span className="ml-2 bg-green-100 text-green-800 px-2 py-1 rounded text-sm font-semibold">
                          {item.finalEstimate}
                        </span>
                      )}
                    </div>
                    <div className="mt-2 text-xs text-gray-500">
                      {Object.keys(item.votes).length} / {Object.keys(session.users).length} voted
                    </div>
                  </div>
                ))}

                {session.items.length === 0 && (
                  <p className="text-gray-500 text-center py-8 text-sm">
                    No items yet. {isHost ? 'Add your first item!' : 'Waiting for host to add items.'}
                  </p>
                )}
              </div>
            </div>

            {/* Participants */}
            <div className="bg-white rounded-lg shadow p-6 mt-6">
              <h2 className="text-xl font-semibold mb-4">Participants ({Object.keys(session.users).length})</h2>
              <div className="space-y-2">
                {Object.values(session.users).map((user) => (
                  <div key={user.id} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                    <span className="text-gray-900">
                      {user.name}
                      {user.isHost && <span className="ml-2 text-xs text-blue-600">(Host)</span>}
                    </span>
                    {currentItem && !currentItem.revealed && currentItem.votes[user.id] && (
                      <span className="text-green-600">✓</span>
                    )}
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Right Column: Voting Area */}
          <div className="lg:col-span-2">
            {currentItem ? (
              <div className="bg-white rounded-lg shadow p-6">
                <div className="mb-6">
                  <h2 className="text-2xl font-bold text-gray-900 mb-2">{currentItem.title}</h2>
                  {currentItem.description && (
                    <p className="text-gray-600">{currentItem.description}</p>
                  )}
                </div>

                {/* Voting Cards */}
                {!currentItem.revealed && (
                  <div>
                    <h3 className="text-lg font-semibold mb-4">Select your estimate:</h3>
                    <div className="grid grid-cols-4 sm:grid-cols-6 gap-3 mb-6">
                      {CARD_VALUES.map((value) => (
                        <button
                          key={value}
                          onClick={() => handleVote(value)}
                          className={`aspect-[2/3] rounded-lg border-2 text-2xl font-bold transition transform hover:scale-105 ${
                            selectedVote === value
                              ? 'border-blue-500 bg-blue-500 text-white shadow-lg'
                              : 'border-gray-300 bg-white text-gray-700 hover:border-blue-300'
                          }`}
                        >
                          {value}
                        </button>
                      ))}
                    </div>
                  </div>
                )}

                {/* Results */}
                {currentItem.revealed && (
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold mb-4">Results:</h3>
                    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
                      {Object.entries(currentItem.votes).map(([userId, vote]) => {
                        const user = session.users[userId];
                        return (
                          <div key={userId} className="bg-gray-50 rounded-lg p-4 text-center">
                            <div className="text-3xl font-bold text-blue-600 mb-2">{vote}</div>
                            <div className="text-sm text-gray-600">{user?.name || 'Unknown'}</div>
                          </div>
                        );
                      })}
                    </div>

                    {isHost && !currentItem.finalEstimate && (
                      <div className="mt-6">
                        <h4 className="text-sm font-semibold mb-2">Set Final Estimate:</h4>
                        <div className="flex gap-2 flex-wrap">
                          {CARD_VALUES.map((value) => (
                            <button
                              key={value}
                              onClick={() => handleSetFinalEstimate(value)}
                              className="px-4 py-2 border-2 border-gray-300 rounded-lg hover:border-green-500 hover:bg-green-50 transition"
                            >
                              {value}
                            </button>
                          ))}
                        </div>
                      </div>
                    )}

                    {currentItem.finalEstimate && (
                      <div className="mt-6 p-4 bg-green-50 border-2 border-green-500 rounded-lg">
                        <p className="text-center text-lg">
                          <span className="font-semibold">Final Estimate:</span>{' '}
                          <span className="text-2xl font-bold text-green-600">{currentItem.finalEstimate}</span>
                        </p>
                      </div>
                    )}
                  </div>
                )}

                {/* Host Controls */}
                {isHost && (
                  <div className="flex gap-3 pt-4 border-t border-gray-200">
                    {!currentItem.revealed ? (
                      <button
                        onClick={handleRevealVotes}
                        className="flex-1 bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition font-semibold"
                      >
                        Reveal Votes
                      </button>
                    ) : (
                      <button
                        onClick={handleResetVotes}
                        className="flex-1 bg-orange-600 text-white px-6 py-3 rounded-lg hover:bg-orange-700 transition font-semibold"
                      >
                        Reset Votes
                      </button>
                    )}
                  </div>
                )}
              </div>
            ) : (
              <div className="bg-white rounded-lg shadow p-12 text-center">
                <div className="text-6xl mb-4">🃏</div>
                <h2 className="text-2xl font-bold text-gray-900 mb-2">No Item Selected</h2>
                <p className="text-gray-600">
                  {isHost
                    ? 'Select an item from the list to start voting'
                    : 'Waiting for the host to select an item'}
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
