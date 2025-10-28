import { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import { createSession } from '@/lib/api';

export default function Home() {
  const [name, setName] = useState('');
  const [sessionName, setSessionName] = useState('');
  const [isCreating, setIsCreating] = useState(false);
  const [isJoining, setIsJoining] = useState(false);
  const [error, setError] = useState('');
  const router = useRouter();
  const { join: joinSessionId, error: urlError } = router.query;

  // Display error from URL parameter if present
  useEffect(() => {
    if (urlError && typeof urlError === 'string') {
      setError(decodeURIComponent(urlError));
      setIsJoining(false); // Reset joining state
      // Clear error from URL but keep the join parameter
      const params = new URLSearchParams(window.location.search);
      params.delete('error');
      const newUrl = params.toString() ? `${router.pathname}?${params.toString()}` : router.pathname;
      router.replace(newUrl, undefined, { shallow: true });
    }
  }, [urlError, router]);

  const handleCreateSession = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!name.trim() || !sessionName.trim()) {
      setError('Please enter both your name and session name');
      return;
    }

    setIsCreating(true);

    try {
      const { sessionId, hostId } = await createSession(sessionName, name);
      router.push(`/session/${sessionId}?userId=${hostId}&userName=${encodeURIComponent(name)}`);
    } catch (err) {
      setError('Failed to create session. Please try again.');
      setIsCreating(false);
    }
  };

  const handleJoinSession = (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!name.trim()) {
      setError('Please enter your name');
      return;
    }

    setIsJoining(true);

    // If there's a join session ID in the URL, use it
    if (joinSessionId) {
      router.push(`/session/${joinSessionId}?userName=${encodeURIComponent(name)}`);
      return;
    }

    // Otherwise, prompt for session ID
    const sessionId = prompt('Enter Session ID:');
    if (sessionId) {
      router.push(`/session/${sessionId}?userName=${encodeURIComponent(name)}`);
    } else {
      setIsJoining(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <div className="max-w-md w-full bg-white rounded-2xl shadow-xl p-8">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">🃏 Poker Planning</h1>
          <p className="text-gray-600">Collaborate and estimate together</p>
          {joinSessionId && (
            <p className="text-sm text-green-600 mt-2 font-semibold">
              Ready to join session! Enter your name below.
            </p>
          )}
        </div>

        <form className="space-y-6">
          <div>
            <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-2">
              Your Name
            </label>
            <input
              type="text"
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
              placeholder="Enter your name"
              required
              autoFocus={!!joinSessionId}
            />
          </div>

          {joinSessionId ? (
            <div className="border-t border-gray-200 pt-6">
              <button
                type="button"
                onClick={handleJoinSession}
                disabled={isJoining}
                className="w-full bg-blue-600 text-white font-semibold py-3 px-6 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isJoining ? 'Joining...' : 'Join Session'}
              </button>
            </div>
          ) : (
            <>
              <div className="border-t border-gray-200 pt-6">
                <h2 className="text-lg font-semibold text-gray-900 mb-4">Create New Session</h2>
                <div className="mb-4">
                  <label htmlFor="sessionName" className="block text-sm font-medium text-gray-700 mb-2">
                    Session Name
                  </label>
                  <input
                    type="text"
                    id="sessionName"
                    value={sessionName}
                    onChange={(e) => setSessionName(e.target.value)}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
                    placeholder="Sprint Planning - Dec 2023"
                  />
                </div>
                <button
                  type="submit"
                  onClick={handleCreateSession}
                  disabled={isCreating}
                  className="w-full bg-blue-600 text-white font-semibold py-3 px-6 rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {isCreating ? 'Creating...' : 'Create Session'}
                </button>
              </div>

              <div className="border-t border-gray-200 pt-6">
                <button
                  type="button"
                  onClick={handleJoinSession}
                  className="w-full bg-white text-blue-600 font-semibold py-3 px-6 rounded-lg border-2 border-blue-600 hover:bg-blue-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition"
                >
                  Join Existing Session
                </button>
              </div>
            </>
          )}

          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
              {error}
            </div>
          )}
        </form>
      </div>
    </div>
  );
}
