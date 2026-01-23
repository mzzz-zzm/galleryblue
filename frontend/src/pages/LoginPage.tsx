import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useMutation } from '@connectrpc/connect-query';
import { login as loginRpc } from '../gen/users/v1/user-AuthService_connectquery';
import { useAuth } from '../context/AuthContext';
import { transport } from '../lib/transport';

export const LoginPage = () => {
    const navigate = useNavigate();
    const { login } = useAuth();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');

    const loginMutation = useMutation(loginRpc, { transport });

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        try {
            const response = await loginMutation.mutateAsync({
                email,
                password,
            });

            login(response.sessionToken, {
                userId: response.userId,
                displayName: response.displayName,
                email: response.email,
            });

            navigate('/');
        } catch (err) {
            setError((err as Error).message || 'Login failed');
        }
    };

    return (
        <div className="auth-page">
            <div className="auth-container">
                <h1 className="auth-title">Welcome Back</h1>

                <form onSubmit={handleSubmit} className="auth-form">
                    {error && <div className="error-message">{error}</div>}

                    <div className="form-group">
                        <label htmlFor="email">Email</label>
                        <input
                            id="email"
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                            placeholder="you@example.com"
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="password">Password</label>
                        <input
                            id="password"
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                            placeholder="••••••••"
                        />
                    </div>

                    <button
                        type="submit"
                        className="btn btn-primary btn-full"
                        disabled={loginMutation.isPending}
                    >
                        {loginMutation.isPending ? 'Logging in...' : 'Login'}
                    </button>
                </form>

                <p className="auth-switch">
                    Don't have an account? <Link to="/register">Register</Link>
                </p>
            </div>
        </div>
    );
};
