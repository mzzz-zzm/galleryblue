import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useMutation } from '@connectrpc/connect-query';
import { register } from '../gen/users/v1/user-AuthService_connectquery';
import { useAuth } from '../context/AuthContext';
import { transport } from '../lib/transport';

export const RegisterPage = () => {
    const navigate = useNavigate();
    const { login } = useAuth();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [displayName, setDisplayName] = useState('');
    const [error, setError] = useState('');

    const registerMutation = useMutation(register, { transport });

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        try {
            const response = await registerMutation.mutateAsync({
                email,
                password,
                displayName,
            });

            // Auto-login after registration
            login('registered', {
                userId: response.userId,
                displayName: response.displayName,
                email: response.email,
            });

            navigate('/');
        } catch (err) {
            setError((err as Error).message || 'Registration failed');
        }
    };

    return (
        <div className="auth-page">
            <div className="auth-container">
                <h1 className="auth-title">Create Account</h1>

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
                        <label htmlFor="displayName">Display Name</label>
                        <input
                            id="displayName"
                            type="text"
                            value={displayName}
                            onChange={(e) => setDisplayName(e.target.value)}
                            placeholder="Your name"
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
                        disabled={registerMutation.isPending}
                    >
                        {registerMutation.isPending ? 'Creating Account...' : 'Register'}
                    </button>
                </form>

                <p className="auth-switch">
                    Already have an account? <Link to="/login">Login</Link>
                </p>
            </div>
        </div>
    );
};
