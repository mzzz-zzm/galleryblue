import React, { useState, useMemo } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useMutation } from '@connectrpc/connect-query';
import { uploadImage } from '../gen/users/v1/user-ImageService_connectquery';
import { useAuth } from '../context/AuthContext';
import { createAuthenticatedTransport } from '../lib/transport';

export const UploadPage = () => {
    const navigate = useNavigate();
    const { user, isAuthenticated } = useAuth();
    const [file, setFile] = useState<File | null>(null);
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [error, setError] = useState('');
    const [preview, setPreview] = useState<string | null>(null);

    const authTransport = useMemo(() => {
        return user?.userId ? createAuthenticatedTransport(user.userId) : null;
    }, [user?.userId]);

    const uploadMutation = useMutation(uploadImage, {
        transport: authTransport ?? undefined,
    });

    if (!isAuthenticated || !user) {
        navigate('/login');
        return null;
    }

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const selectedFile = e.target.files?.[0];
        if (selectedFile) {
            if (selectedFile.type !== 'image/jpeg') {
                setError('Only JPEG images are supported');
                return;
            }
            if (selectedFile.size > 5 * 1024 * 1024) {
                setError('Image too large (max 5MB)');
                return;
            }
            setFile(selectedFile);
            setError('');
            // Create preview
            const reader = new FileReader();
            reader.onload = () => setPreview(reader.result as string);
            reader.readAsDataURL(selectedFile);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (!file) {
            setError('Please select an image');
            return;
        }

        try {
            const arrayBuffer = await file.arrayBuffer();
            const data = new Uint8Array(arrayBuffer);

            await uploadMutation.mutateAsync({
                filename: file.name,
                contentType: 'image/jpeg',
                data,
                title,
                description,
            });

            navigate('/my-images');
        } catch (err) {
            setError((err as Error).message || 'Upload failed');
        }
    };

    return (
        <div className="auth-page">
            <div className="auth-container" style={{ maxWidth: '500px' }}>
                <h1 className="auth-title">Upload Image</h1>

                <form onSubmit={handleSubmit} className="auth-form">
                    {error && <div className="error-message">{error}</div>}

                    <div className="form-group">
                        <label htmlFor="file">Select JPEG Image</label>
                        <input
                            id="file"
                            type="file"
                            accept="image/jpeg"
                            onChange={handleFileChange}
                            required
                        />
                    </div>

                    {preview && (
                        <div style={{ marginBottom: '1rem' }}>
                            <img
                                src={preview}
                                alt="Preview"
                                style={{ maxWidth: '100%', maxHeight: '200px', borderRadius: '8px' }}
                            />
                        </div>
                    )}

                    <div className="form-group">
                        <label htmlFor="title">Title (optional)</label>
                        <input
                            id="title"
                            type="text"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            placeholder="Give your image a title"
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="description">Description (optional)</label>
                        <textarea
                            id="description"
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            placeholder="Describe your image"
                            rows={3}
                            style={{ width: '100%', padding: '0.5rem', borderRadius: '4px', border: '1px solid #ccc' }}
                        />
                    </div>

                    <button
                        type="submit"
                        className="btn btn-primary btn-full"
                        disabled={uploadMutation.isPending || !file}
                    >
                        {uploadMutation.isPending ? 'Uploading...' : 'Upload Image'}
                    </button>
                </form>

                <p className="auth-switch">
                    <Link to="/my-images">← My Images</Link> | <Link to="/">Home</Link>
                </p>
            </div>
        </div>
    );
};
