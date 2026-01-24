import { useMemo, useState } from 'react';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import { deleteImage, getImage } from '../gen/users/v1/user-ImageService_connectquery';
import { useAuth } from '../context/AuthContext';
import { transport, createAuthenticatedTransport } from '../lib/transport';

interface ImageCardProps {
    id: string;
    title: string;
    ownerName: string;
    createdAt: string;
    isOwner: boolean;
    thumbnail?: Uint8Array;
    onDelete?: () => void;
}

export const ImageCard: React.FC<ImageCardProps> = ({
    id,
    title,
    ownerName,
    createdAt,
    isOwner,
    thumbnail,
    onDelete,
}) => {
    const { user } = useAuth();
    const [confirmDelete, setConfirmDelete] = useState(false);
    const [showFullImage, setShowFullImage] = useState(false);

    const authTransport = useMemo(() => {
        return user?.userId ? createAuthenticatedTransport(user.userId) : null;
    }, [user?.userId]);

    // Fetch full image only when modal is open
    const { data: fullImageData, isLoading: loadingFullImage } = useQuery(
        getImage,
        { id },
        { transport, enabled: showFullImage }
    );

    const deleteMutation = useMutation(deleteImage, {
        transport: authTransport ?? undefined,
    });

    const handleDelete = async () => {
        try {
            await deleteMutation.mutateAsync({ id });
            onDelete?.();
        } catch (err) {
            alert((err as Error).message);
        }
    };

    const formatDate = (dateStr: string) => {
        try {
            return new Date(dateStr).toLocaleDateString();
        } catch {
            return dateStr;
        }
    };

    // Convert bytes to base64 data URL
    const bytesToDataUrl = (bytes: Uint8Array | undefined, contentType = 'image/jpeg') => {
        if (!bytes || bytes.length === 0) return null;
        try {
            let binary = '';
            const chunkSize = 8192;
            for (let i = 0; i < bytes.length; i += chunkSize) {
                const chunk = bytes.slice(i, i + chunkSize);
                binary += String.fromCharCode.apply(null, Array.from(chunk));
            }
            return `data:${contentType};base64,${btoa(binary)}`;
        } catch (e) {
            console.error('Failed to convert image:', e);
            return null;
        }
    };

    const thumbnailUrl = useMemo(() => bytesToDataUrl(thumbnail), [thumbnail]);
    const fullImageUrl = useMemo(
        () => bytesToDataUrl(fullImageData?.data, fullImageData?.contentType),
        [fullImageData]
    );

    return (
        <>
            <div style={{
                border: '1px solid #ddd',
                borderRadius: '8px',
                overflow: 'hidden',
                backgroundColor: '#fff',
                boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
            }}>
                {/* Clickable thumbnail */}
                <div
                    style={{
                        height: '180px',
                        backgroundColor: '#f0f0f0',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        cursor: 'pointer',
                    }}
                    onClick={() => setShowFullImage(true)}
                    title="Click to view full size"
                >
                    {thumbnailUrl ? (
                        <img
                            src={thumbnailUrl}
                            alt={title}
                            style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                        />
                    ) : (
                        <span style={{ color: '#999' }}>No preview</span>
                    )}
                </div>

                {/* Info section */}
                <div style={{ padding: '1rem' }}>
                    <h3 style={{ margin: '0 0 0.5rem', fontSize: '1rem' }}>{title}</h3>
                    <p style={{ margin: '0 0 0.25rem', fontSize: '0.85rem', color: '#666' }}>
                        By: {ownerName}
                    </p>
                    <p style={{ margin: '0', fontSize: '0.75rem', color: '#999' }}>
                        {formatDate(createdAt)}
                    </p>

                    {isOwner && (
                        <div style={{ marginTop: '1rem', display: 'flex', gap: '0.5rem' }}>
                            {!confirmDelete ? (
                                <button
                                    onClick={() => setConfirmDelete(true)}
                                    style={{
                                        padding: '0.25rem 0.5rem',
                                        fontSize: '0.85rem',
                                        background: '#ff4444',
                                        color: 'white',
                                        border: 'none',
                                        borderRadius: '4px',
                                        cursor: 'pointer',
                                    }}
                                >
                                    Delete
                                </button>
                            ) : (
                                <>
                                    <button
                                        onClick={handleDelete}
                                        disabled={deleteMutation.isPending}
                                        style={{
                                            padding: '0.25rem 0.5rem',
                                            fontSize: '0.85rem',
                                            background: '#ff0000',
                                            color: 'white',
                                            border: 'none',
                                            borderRadius: '4px',
                                            cursor: 'pointer',
                                        }}
                                    >
                                        {deleteMutation.isPending ? '...' : 'Confirm'}
                                    </button>
                                    <button
                                        onClick={() => setConfirmDelete(false)}
                                        style={{
                                            padding: '0.25rem 0.5rem',
                                            fontSize: '0.85rem',
                                            background: '#ccc',
                                            border: 'none',
                                            borderRadius: '4px',
                                            cursor: 'pointer',
                                        }}
                                    >
                                        Cancel
                                    </button>
                                </>
                            )}
                        </div>
                    )}
                </div>
            </div>

            {/* Full-size image modal */}
            {showFullImage && (
                <div
                    style={{
                        position: 'fixed',
                        top: 0,
                        left: 0,
                        right: 0,
                        bottom: 0,
                        backgroundColor: 'rgba(0, 0, 0, 0.9)',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        zIndex: 1000,
                        cursor: 'pointer',
                    }}
                    onClick={() => setShowFullImage(false)}
                >
                    {loadingFullImage ? (
                        <div style={{ color: 'white', fontSize: '1.5rem' }}>Loading...</div>
                    ) : fullImageUrl ? (
                        <div style={{ position: 'relative', maxWidth: '90vw', maxHeight: '90vh' }}>
                            <img
                                src={fullImageUrl}
                                alt={title}
                                style={{
                                    maxWidth: '90vw',
                                    maxHeight: '90vh',
                                    objectFit: 'contain',
                                    borderRadius: '4px',
                                }}
                                onClick={(e) => e.stopPropagation()}
                            />
                            <div style={{
                                position: 'absolute',
                                bottom: '-40px',
                                left: 0,
                                right: 0,
                                textAlign: 'center',
                                color: 'white',
                            }}>
                                <strong>{title}</strong> by {ownerName}
                            </div>
                            <button
                                onClick={() => setShowFullImage(false)}
                                style={{
                                    position: 'absolute',
                                    top: '-40px',
                                    right: 0,
                                    background: 'transparent',
                                    border: 'none',
                                    color: 'white',
                                    fontSize: '2rem',
                                    cursor: 'pointer',
                                }}
                            >
                                ✕
                            </button>
                        </div>
                    ) : (
                        <div style={{ color: 'white' }}>Failed to load image</div>
                    )}
                </div>
            )}
        </>
    );
};
