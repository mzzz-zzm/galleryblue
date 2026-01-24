import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { TransportProvider } from "@connectrpc/connect-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { transport } from "./lib/transport";
import { AuthProvider } from "./context/AuthContext";
import { Layout } from "./components/Layout";
import { HomePage } from "./pages/HomePage";
import { RegisterPage } from "./pages/RegisterPage";
import { LoginPage } from "./pages/LoginPage";
import { UpdatePage } from "./pages/UpdatePage";
import { UploadPage } from "./pages/UploadPage";
import { GalleryPage } from "./pages/GalleryPage";
import { MyImagesPage } from "./pages/MyImagesPage";

const queryClient = new QueryClient();

function App() {
    return (
        <BrowserRouter>
            <TransportProvider transport={transport}>
                <QueryClientProvider client={queryClient}>
                    <AuthProvider>
                        <Routes>
                            {/* Auth pages without sidebar */}
                            <Route path="/register" element={<RegisterPage />} />
                            <Route path="/login" element={<LoginPage />} />

                            {/* Main pages with sidebar */}
                            <Route path="/" element={<Layout><HomePage /></Layout>} />
                            <Route path="/update" element={<Layout><UpdatePage /></Layout>} />
                            <Route path="/upload" element={<Layout><UploadPage /></Layout>} />
                            <Route path="/gallery" element={<Layout><GalleryPage /></Layout>} />
                            <Route path="/my-images" element={<Layout><MyImagesPage /></Layout>} />
                        </Routes>
                    </AuthProvider>
                </QueryClientProvider>
            </TransportProvider>
        </BrowserRouter>
    );
}

export default App;

