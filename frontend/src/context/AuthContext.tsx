import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';

interface User {
    userId: string;
    displayName: string;
    email: string;
}

interface AuthContextType {
    user: User | null;
    sessionToken: string | null;
    login: (token: string, user: User) => void;
    logout: () => void;
    updateUser: (user: User) => void;
    isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [sessionToken, setSessionToken] = useState<string | null>(null);

    // Load from localStorage on mount
    useEffect(() => {
        const storedToken = localStorage.getItem('sessionToken');
        const storedUser = localStorage.getItem('user');
        if (storedToken && storedUser) {
            setSessionToken(storedToken);
            setUser(JSON.parse(storedUser));
        }
    }, []);

    const login = (token: string, userData: User) => {
        setSessionToken(token);
        setUser(userData);
        localStorage.setItem('sessionToken', token);
        localStorage.setItem('user', JSON.stringify(userData));
    };

    const logout = () => {
        setSessionToken(null);
        setUser(null);
        localStorage.removeItem('sessionToken');
        localStorage.removeItem('user');
    };

    const updateUser = (userData: User) => {
        setUser(userData);
        localStorage.setItem('user', JSON.stringify(userData));
    };

    return (
        <AuthContext.Provider
            value={{
                user,
                sessionToken,
                login,
                logout,
                updateUser,
                isAuthenticated: !!sessionToken,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
};
