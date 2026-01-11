import { createConnectTransport } from "@connectrpc/connect-web";

// This transport communicates with the Go backend
export const transport = createConnectTransport({
    baseUrl: "/api",
});
