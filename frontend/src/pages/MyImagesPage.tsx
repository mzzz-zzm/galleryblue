import { useMemo } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useQuery } from '@connectrpc/connect-query';
import { listMyImages } from '../gen/users/v1/user-ImageService_connectquery';
import { useAuth } from '../context/AuthContext';
import { createAuthenticatedTransport } from '../lib/transport';
import { ImageCard } from '../components/ImageCard';

export const MyImagesPage = () => {
    const navigate = useNavigate();
    const { user, isAuthenticated } = useAuth();

    const authTransport = useMemo(() => {
        return user?.userId ? createAuthenticatedTransport(user.userId) : null;
    }, [user?.userId]);

    const { data, isLoading, error, refetch } = useQuery(
        listMyImages,
        { limit: 50, offset: 0 },
        { transport: authTransport ?? undefined, enabled: !!authTransport }
    );

    if (!isAuthenticated || !user) {
        navigate('/login');
        return null;
    }

    if (isLoading) {
        return (
            <div className="auth-page">
                <div className="auth-container">
                    <p>Loading your images...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="auth-page">
                <div className="auth-container">
                    <div className="error-message">{(error as Error).message}</div>
                </div>
            </div>
        );
    }

    const images = data?.images || [];

    return (
        <div className="auth-page">
            <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '2rem' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
                    <h1>My Images</h1>
                    <div>
                        <Link to="/upload" className="btn btn-primary" style={{ marginRight: '1rem' }}>
                            Upload New
                        </Link>
                        <Link to="/gallery" style={{ marginRight: '1rem' }}>Gallery</Link>
                        <Link to="/">Home</Link>
                    </div>
                </div>

                {images.length === 0 ? (
                    <p>You haven't uploaded any images yet. <Link to="/upload">Upload your first image!</Link></p>
                ) : (
                    <div style={{
                        display: 'grid',
                        gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))',
                        gap: '1.5rem'
                    }}>
                        {images.map((img) => (
                            <ImageCard
                                key={img.id}
                                id={img.id}
                                title={img.title || img.filename}
                                ownerName={img.ownerDisplayName}
                                createdAt={img.createdAt}
                                isOwner={true}
                                thumbnail={img.thumbnail}
                                onDelete={() => refetch()}
                            />
                        ))}
                    </div>
                )}

                <p style={{ marginTop: '2rem', color: '#666' }}>
                    Total: {data?.total || 0} images
                </p>
            </div>
        </div>
    );
};
