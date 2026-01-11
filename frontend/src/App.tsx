import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { TransportProvider } from "@connectrpc/connect-query";
import { transport } from "./lib/transport";
import { UserProfile } from "./components/UserProfile";

const queryClient = new QueryClient();

function App() {
    return (
        <TransportProvider transport={transport}>
            <QueryClientProvider client={queryClient}>
                <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center">
                    <UserProfile userId="123" />
                </div>
            </QueryClientProvider>
        </TransportProvider>
    );
}

export default App;
