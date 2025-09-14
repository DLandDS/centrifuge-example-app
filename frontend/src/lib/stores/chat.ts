import { writable } from 'svelte/store';
import { Centrifuge } from 'centrifuge';
import type { Message } from '$lib/types';

interface ChatState {
    isConnected: boolean;
    messages: Message[];
    activeTopic: string;
    topics: string[];
}

const initialState: ChatState = {
    isConnected: false,
    messages: [],
    activeTopic: '',
    topics: ['general', 'tech', 'random']
};

export const chatStore = writable<ChatState>(initialState);

let centrifuge: Centrifuge | null = null;
let currentSubscription: any = null;

export const connectToCentrifuge = (token: string) => {
    if (centrifuge) {
        centrifuge.disconnect();
    }

    centrifuge = new Centrifuge('ws://localhost:8000/connection/websocket', {
        token: token
    });

    centrifuge.on('connecting', () => {
        console.log('Connecting to Centrifuge...');
    });

    centrifuge.on('connected', () => {
        console.log('Connected to Centrifuge');
        chatStore.update(state => ({ ...state, isConnected: true }));
    });

    centrifuge.on('disconnected', () => {
        console.log('Disconnected from Centrifuge');
        chatStore.update(state => ({ ...state, isConnected: false }));
    });

    centrifuge.connect();
};

export const subscribeToTopic = (topic: string) => {
    if (!centrifuge) return;

    // Unsubscribe from current topic
    if (currentSubscription) {
        currentSubscription.unsubscribe();
        currentSubscription.removeAllListeners();
        currentSubscription = null;
    }

    // Subscribe to new topic
    const channelName = `chat:${topic}`;
    currentSubscription = centrifuge.newSubscription(channelName);

    currentSubscription.on('publication', (ctx: any) => {
        const message: Message = {
            id: `${Date.now()}-${Math.random()}`,
            username: ctx.data.username,
            content: ctx.data.content,
            timestamp: ctx.data.timestamp || Date.now(),
            topic: topic
        };

        chatStore.update(state => ({
            ...state,
            messages: [...state.messages, message]
        }));
    });

    currentSubscription.on('subscribed', () => {
        console.log(`Subscribed to ${channelName}`);
        chatStore.update(state => ({ 
            ...state, 
            activeTopic: topic,
            messages: [] // Clear messages when switching topics
        }));
    });

    currentSubscription.subscribe();
};

export const sendMessage = (content: string, username: string) => {
    if (!currentSubscription) return;

    const messageData = {
        username,
        content,
        timestamp: Date.now()
    };

    currentSubscription.publish(messageData);
};

export const disconnect = () => {
    if (currentSubscription) {
        currentSubscription.unsubscribe();
        currentSubscription = null;
    }
    
    if (centrifuge) {
        centrifuge.disconnect();
        centrifuge = null;
    }
    
    chatStore.set(initialState);
};