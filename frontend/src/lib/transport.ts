import { createConnectTransport } from "@connectrpc/connect-web";

// This transport communicates with the Go backend
export const transport = createConnectTransport({
    baseUrl: "/api",
});

// Create an authenticated transport with user ID header
export const createAuthenticatedTransport = (userId: string) => {
    return createConnectTransport({
        baseUrl: "/api",
        interceptors: [
            (next) => async (req) => {
                req.header.set("X-User-ID", userId);
                return next(req);
            },
        ],
    });
};
