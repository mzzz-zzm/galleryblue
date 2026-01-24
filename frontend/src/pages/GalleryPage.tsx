import { Link } from 'react-router-dom';
import { useQuery } from '@connectrpc/connect-query';
import { listImages } from '../gen/users/v1/user-ImageService_connectquery';
import { transport } from '../lib/transport';
import { ImageCard } from '../components/ImageCard';
import { useAuth } from '../context/AuthContext';

export const GalleryPage = () => {
    const { isAuthenticated } = useAuth();

    const { data, isLoading, error } = useQuery(listImages, { limit: 50, offset: 0 }, { transport });

    if (isLoading) {
        return (
            <div className="auth-page">
                <div className="auth-container">
                    <p>Loading gallery...</p>
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
                    <h1>Gallery</h1>
                    <div>
                        {isAuthenticated && (
                            <>
                                <Link to="/upload" className="btn btn-primary" style={{ marginRight: '1rem' }}>
                                    Upload Image
                                </Link>
                                <Link to="/my-images" className="btn" style={{ marginRight: '1rem' }}>
                                    My Images
                                </Link>
                            </>
                        )}
                        <Link to="/">Home</Link>
                    </div>
                </div>

                {images.length === 0 ? (
                    <p>No images yet. {isAuthenticated ? <Link to="/upload">Upload the first one!</Link> : 'Login to upload images.'}</p>
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
                                isOwner={false}
                                thumbnail={img.thumbnail}
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
