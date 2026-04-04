import {useState} from "react";
import {apiRequest, RequestMethods} from "@/api/api";

export default function usePost() {
    const [result, setResult] = useState<string | null>(null)
    const postData = async (endpoint: string, body?: Record<string, any>, method?: RequestMethods) => {
        try {
            await apiRequest(endpoint, {method: method || "POST", body})
            setResult(null)
        } catch (error) {
            if (error instanceof Error) {
                setResult(error.message);
            } else {
                setResult("Error Occurred")
            }
        }
    }

    const clearResult = () => {
        setResult(null)
    }

    return {postData, result, clearResult}
}
