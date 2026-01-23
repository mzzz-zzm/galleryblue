import { useState, useMemo } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useMutation } from '@connectrpc/connect-query';
import { updateUser } from '../gen/users/v1/user-UserService_connectquery';
import { useAuth } from '../context/AuthContext';
import { createAuthenticatedTransport } from '../lib/transport';

export const UpdatePage = () => {
    const navigate = useNavigate();
    const { user, isAuthenticated, updateUser: updateUserContext } = useAuth();
    const [currentPassword, setCurrentPassword] = useState('');
    const [newDisplayName, setNewDisplayName] = useState('');
    const [newEmail, setNewEmail] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    // Create authenticated transport with user ID
    const authTransport = useMemo(() => {
        return user?.userId ? createAuthenticatedTransport(user.userId) : null;
    }, [user?.userId]);

    const updateMutation = useMutation(updateUser, {
        transport: authTransport ?? undefined,
        onSuccess: () => { },
    });

    // Redirect if not authenticated
    if (!isAuthenticated || !user) {
        navigate('/login');
        return null;
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        if (!currentPassword) {
            setError('Current password is required');
            return;
        }

        try {
            const response = await updateMutation.mutateAsync({
                currentPassword,
                newDisplayName: newDisplayName || undefined,
                newEmail: newEmail || undefined,
                newPassword: newPassword || undefined,
            });

            updateUserContext({
                userId: response.userId,
                displayName: response.displayName,
                email: response.email,
            });

            setSuccess('Profile updated successfully!');
            setCurrentPassword('');
            setNewPassword('');
        } catch (err) {
            setError((err as Error).message || 'Update failed');
        }
    };

    return (
        <div className="auth-page">
            <div className="auth-container">
                <h1 className="auth-title">Update Profile</h1>

                <form onSubmit={handleSubmit} className="auth-form">
                    {error && <div className="error-message">{error}</div>}
                    {success && <div className="success-message">{success}</div>}

                    <div className="form-group">
                        <label htmlFor="currentPassword">Current Password *</label>
                        <input
                            id="currentPassword"
                            type="password"
                            value={currentPassword}
                            onChange={(e) => setCurrentPassword(e.target.value)}
                            required
                            placeholder="Enter current password to confirm"
                        />
                    </div>

                    <hr className="form-divider" />

                    <div className="form-group">
                        <label htmlFor="newDisplayName">New Display Name</label>
                        <input
                            id="newDisplayName"
                            type="text"
                            value={newDisplayName}
                            onChange={(e) => setNewDisplayName(e.target.value)}
                            placeholder={user.displayName || 'Leave blank to keep current'}
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="newEmail">New Email</label>
                        <input
                            id="newEmail"
                            type="email"
                            value={newEmail}
                            onChange={(e) => setNewEmail(e.target.value)}
                            placeholder={user.email || 'Leave blank to keep current'}
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="newPassword">New Password</label>
                        <input
                            id="newPassword"
                            type="password"
                            value={newPassword}
                            onChange={(e) => setNewPassword(e.target.value)}
                            placeholder="Leave blank to keep current"
                        />
                    </div>

                    <button
                        type="submit"
                        className="btn btn-primary btn-full"
                        disabled={updateMutation.isPending}
                    >
                        {updateMutation.isPending ? 'Updating...' : 'Update Profile'}
                    </button>
                </form>

                <p className="auth-switch">
                    <Link to="/">← Back to Home</Link>
                </p>
            </div>
        </div>
    );
};
