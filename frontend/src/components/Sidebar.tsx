import { NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export const Sidebar = () => {
    const { user, isAuthenticated, logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/');
    };

    return (
        <aside className="sidebar">
            <div className="sidebar-header">
                <h2 className="sidebar-logo">GalleryBlue</h2>
            </div>

            {isAuthenticated && user && (
                <div className="sidebar-user">
                    <div className="sidebar-avatar">
                        {(user.displayName || user.email).charAt(0).toUpperCase()}
                    </div>
                    <span className="sidebar-username">{user.displayName || 'User'}</span>
                </div>
            )}

            <nav className="sidebar-nav">
                <NavLink to="/" className={({ isActive }) => `sidebar-link ${isActive ? 'active' : ''}`}>
                    <span className="sidebar-icon">🏠</span>
                    Home
                </NavLink>

                <NavLink to="/gallery" className={({ isActive }) => `sidebar-link ${isActive ? 'active' : ''}`}>
                    <span className="sidebar-icon">🖼️</span>
                    Gallery
                </NavLink>

                {isAuthenticated ? (
                    <NavLink to="/upload" className={({ isActive }) => `sidebar-link ${isActive ? 'active' : ''}`}>
                        <span className="sidebar-icon">⬆️</span>
                        Upload Image
                    </NavLink>
                ) : (
                    <span className="sidebar-link disabled" title="Login to upload images">
                        <span className="sidebar-icon">⬆️</span>
                        Upload Image
                        <span className="sidebar-lock">🔒</span>
                    </span>
                )}

                {isAuthenticated && (
                    <>
                        <NavLink to="/my-images" className={({ isActive }) => `sidebar-link ${isActive ? 'active' : ''}`}>
                            <span className="sidebar-icon">📂</span>
                            My Images
                        </NavLink>

                        <NavLink to="/update" className={({ isActive }) => `sidebar-link ${isActive ? 'active' : ''}`}>
                            <span className="sidebar-icon">👤</span>
                            Profile Settings
                        </NavLink>
                    </>
                )}
            </nav>

            <div className="sidebar-footer">
                {isAuthenticated ? (
                    <button onClick={handleLogout} className="sidebar-link sidebar-logout">
                        <span className="sidebar-icon">🚪</span>
                        Logout
                    </button>
                ) : (
                    <NavLink to="/login" className="sidebar-link">
                        <span className="sidebar-icon">🔑</span>
                        Login
                    </NavLink>
                )}
            </div>
        </aside>
    );
};
