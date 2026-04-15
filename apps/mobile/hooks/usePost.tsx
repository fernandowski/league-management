import {useCallback, useState} from "react";
import {apiRequest, RequestMethods} from "@/api/api";

interface PostResult {
    ok: boolean
    error: string | null
}

export default function usePost() {
    const [result, setResult] = useState<string | null>(null)
    const postData = useCallback(async (endpoint: string, body?: Record<string, any>, method?: RequestMethods) => {
        try {
            await apiRequest(endpoint, {method: method || "POST", body})
            setResult(null)
            return {ok: true, error: null} satisfies PostResult;
        } catch (error) {
            let message = "Error Occurred";
            if (error instanceof Error) {
                message = error.message;
            }
            setResult(message)
            return {ok: false, error: message} satisfies PostResult;
        }
    }, []);

    const clearResult = useCallback(() => {
        setResult(null)
    }, []);

    return {postData, result, clearResult}
}
