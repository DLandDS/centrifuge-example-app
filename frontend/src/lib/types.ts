export interface User {
    id: string;
    username: string;
    email: string;
}

export interface LoginRequest {
    username: string;
    password: string;
}

export interface LoginResponse {
    token: string;
    centrifuge_token: string;
    user: User;
}

export interface Message {
    id: string;
    username: string;
    content: string;
    timestamp: number;
    topic: string;
}