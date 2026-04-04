import AsyncStorage from "@react-native-async-storage/async-storage";
import {jwtDecode, JwtPayload} from "jwt-decode";

type Token = string | null
export const fetchJWT = async (): Promise<Token> => {
    try {
        const token = await AsyncStorage.getItem('jwt');
        return token;
    } catch (error) {
        return "";
    }
}

export const storeJWT = async (jwt: string): Promise<void> => {
    await AsyncStorage.setItem('jwt', jwt);
}

export const removeJWT = async (): Promise<void> => {
    await AsyncStorage.removeItem("jwt");
}

export const JwtIsExpired = (jwt: string): boolean => {
    if (!jwt) {
        return false;
    }
    const {exp}: JwtPayload = jwtDecode(jwt);

    const now: number = new Date().getMilliseconds() / 1000;

    return (exp !== undefined && now >= exp);

}
