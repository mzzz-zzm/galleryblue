import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { TransportProvider } from "@connectrpc/connect-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { transport } from "./lib/transport";
import { AuthProvider } from "./context/AuthContext";
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
                            <Route path="/" element={<HomePage />} />
                            <Route path="/register" element={<RegisterPage />} />
                            <Route path="/login" element={<LoginPage />} />
                            <Route path="/update" element={<UpdatePage />} />
                            <Route path="/upload" element={<UploadPage />} />
                            <Route path="/gallery" element={<GalleryPage />} />
                            <Route path="/my-images" element={<MyImagesPage />} />
                        </Routes>
                    </AuthProvider>
                </QueryClientProvider>
            </TransportProvider>
        </BrowserRouter>
    );
}

export default App;
