import {fetchJWT, JwtIsExpired, removeJWT} from "@/util/jwt-manager";
import {router} from "expo-router";
import {Platform} from "react-native";

export type RequestMethods = 'GET' | 'POST' | 'PUT' | 'DELETE';
export interface RequestOptions {
    method: RequestMethods
    headers?: { [key: string]: string };
    body?: any;
}


const API_URL = process.env.EXPO_PUBLIC_API_URL ?? (Platform.OS === "web" ? "http://localhost:8080" : undefined);


export async function apiRequest(
    endpoint: string,
    {method, headers, body}: RequestOptions,
    secure: boolean = true
) {
    if (!API_URL) {
        throw new Error("API base URL is not configured. Set EXPO_PUBLIC_API_URL.");
    }

    const url: string = `${API_URL}${endpoint}`;
    const jwt: string = await fetchJWT() as string;

    if (secure && JwtIsExpired(jwt)) {
        await removeJWT();
        router.push("/login")
        throw new Error("jwt is expried")
    }

    const defaultHeaders: Record<string, any> = {
        'Content-Type': 'application/json',
        ...(headers || {}),
    };

    if (secure) {
        defaultHeaders["auth"] = jwt
    }

    const options: RequestInit = {
        method,
        headers: defaultHeaders,
        ...(body && {body: JSON.stringify(body)}),
    };

    try {
        const response :Response = await fetch(url, options);

        if (!response.ok) {
            if (response.status === 401) {
                router.push("/login");
                throw new Error("unauthorized");
            }
            const errorResponse = await response.json();
            throw new Error(errorResponse.message || errorResponse.error || 'Request failed');
        }
        return response.json();
    } catch (error) {
        throw error;
    }
}
