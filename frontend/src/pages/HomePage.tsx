import { useAuth } from '../context/AuthContext';
import { Link } from 'react-router-dom';

export const HomePage = () => {
    const { user, isAuthenticated, logout } = useAuth();

    return (
        <div className="home-page">
            <div className="home-container">
                <h1 className="home-title">GalleryBlue</h1>

                {isAuthenticated && user ? (
                    <div className="user-card">
                        <div className="user-avatar">
                            {(user.displayName || user.email).charAt(0).toUpperCase()}
                        </div>
                        <div className="user-info">
                            <h2 className="user-name">{user.displayName || 'User'}</h2>
                            <p className="user-email">{user.email}</p>
                        </div>
                        <div className="user-actions">
                            <Link to="/update" className="btn btn-secondary">
                                Update Profile
                            </Link>
                            <button onClick={logout} className="btn btn-outline">
                                Logout
                            </button>
                        </div>
                    </div>
                ) : (
                    <div className="guest-card">
                        <p className="guest-message">Welcome to GalleryBlue</p>
                        <div className="guest-actions">
                            <Link to="/login" className="btn btn-primary">
                                Login
                            </Link>
                            <Link to="/register" className="btn btn-secondary">
                                Register
                            </Link>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};
