export class APIError extends Error {
    // cast to error
    status: number;

    constructor(message: string, status: number) {
        super(message);
        this.name = "APIError";
        this.status = status;
    }
}

// Helps avoid stuff like form data being sent as JSON
type RequestOptions = Omit<RequestInit, "body"> & {
    body?: unknown;
};

const API_BASE = (import.meta.env.VITE_API_BASE_URL || "/api/v1").replace(/\/$/, "");

async function parseResponse(response: Response): Promise<unknown> {
    if (response.status === 204) return null; // avoid response.json()

    const contentType = response.headers.get("content-type") || "";
    if (contentType.includes("application/json")) {
        return response.json();
    }

    const text = await response.text();
    return text ? { message: text } : null; // normalize object returns
}

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
    const headers = new Headers(options.headers);
    if (!headers.has("Content-Type")) {
        headers.set("Content-Type", "application/json");
    }

    const response = await fetch(`${API_BASE}${path}`, {
        ...options,
        headers,
        credentials: "include",
        body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
    });

    const payload = await parseResponse(response);

    if (!response.ok) {
        let message =
            payload && typeof payload === "object" && "error" in payload && payload.error ?
                String(payload.error)
            :   `Request failed (${response.status})`;
        throw new APIError(message, response.status);
    }

    return payload as T;
}

// Response Types
export type AuthUser = {
    id: number;
    email: string;
    name: string;
    created_at: string;
    updated_at: string;
};

export const api = {
    auth: {
        signup(input: { name: string; email: string; password: string }) {
            return request<AuthUser>("/auth/signup", { method: "POST", body: input });
        },
        login(input: { email: string; password: string }) {
            return request<AuthUser>("/auth/login", { method: "POST", body: input });
        },
    },
};
