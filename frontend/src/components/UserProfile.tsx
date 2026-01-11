import { useQuery } from "@connectrpc/connect-query";
import { getUser } from "../gen/users/v1/user-UserService_connectquery";
import { useState } from "react";

import { transport } from "../lib/transport";

export const UserProfile = ({ userId }: { userId: string }) => {
    const { data, isLoading, error, fetchStatus, status } = useQuery(getUser, { id: userId }, { transport });
    const [manualResult, setManualResult] = useState<string>("");

    const performManualFetch = async () => {
        try {
            const res = await fetch("/api/users.v1.UserService/GetUser", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ id: userId }),
            });
            const text = await res.text();
            setManualResult(`Status: ${res.status}, Body: ${text}`);
        } catch (e) {
            setManualResult(`Error: ${(e as Error).message}`);
        }
    };

    return (
        <div className="p-6 bg-white rounded-lg shadow-md space-y-4">
            <h1 className="text-2xl font-bold">User Profile</h1>

            <div className="p-4 bg-gray-100 rounded text-sm font-mono">
                <p><strong>Query Status:</strong> {status}</p>
                <p><strong>Fetch Status:</strong> {fetchStatus}</p>
                <p><strong>Is Loading:</strong> {isLoading ? "true" : "false"}</p>
            </div>

            {isLoading && <div className="text-blue-500">Loading...</div>}

            {error && (
                <div className="p-4 bg-red-100 text-red-700 rounded">
                    <strong>Error:</strong> {(error as Error).message}
                    <pre className="text-xs mt-2">{JSON.stringify(error, null, 2)}</pre>
                </div>
            )}

            {data && (
                <div className="space-y-2 border-l-4 border-green-500 pl-4">
                    <p><span className="font-semibold">ID:</span> {data.id}</p>
                    <p><span className="font-semibold">Name:</span> {data.name}</p>
                    <p><span className="font-semibold">Email:</span> {data.email}</p>
                </div>
            )}

            <div className="pt-4 border-t">
                <button
                    onClick={performManualFetch}
                    className="px-4 py-2 bg-gray-800 text-white rounded hover:bg-gray-700"
                >
                    Test Manual Fetch
                </button>
                {manualResult && (
                    <div className="mt-2 p-2 bg-yellow-100 rounded text-xs break-all">
                        {manualResult}
                    </div>
                )}
            </div>
        </div>
    );
};
