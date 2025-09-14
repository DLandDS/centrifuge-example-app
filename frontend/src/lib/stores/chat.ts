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
    topics: ['all', 'general', 'tech', 'random']
};

export const chatStore = writable<ChatState>(initialState);

let centrifuge: Centrifuge | null = null;
let currentSubscription: any = null;
let allTopicSubscriptions: Map<string, any> = new Map();

export const connectToCentrifuge = (token: string) => {
    if (centrifuge) {
        centrifuge.disconnect();
    }

    console.log('Creating Centrifuge connection with token:', token.substring(0, 20) + '...');
    
    // Try different connection configuration
    centrifuge = new Centrifuge('ws://localhost:3001/connection/websocket', {
        token: token,
        debug: true,
        name: 'js',
        version: '5.4.0'
    });

    centrifuge.on('connecting', (ctx) => {
        console.log('Connecting to Centrifuge...', ctx);
    });

    centrifuge.on('connected', (ctx) => {
        console.log('Connected to Centrifuge successfully!', ctx);
        chatStore.update(state => ({ ...state, isConnected: true }));
    });

    centrifuge.on('disconnected', (ctx) => {
        console.log('Disconnected from Centrifuge:', ctx);
        chatStore.update(state => ({ ...state, isConnected: false }));
    });

    centrifuge.on('error', (ctx) => {
        console.error('Centrifuge error:', ctx);
    });

    try {
        centrifuge.connect();
    } catch (error) {
        console.error('Failed to connect:', error);
    }
};

export const subscribeToTopic = (topic: string) => {
    if (!centrifuge) return;

    // Clean up all existing subscriptions more thoroughly
    if (currentSubscription) {
        currentSubscription.unsubscribe();
        currentSubscription.removeAllListeners();
        currentSubscription = null;
    }

    // Unsubscribe from all topic subscriptions
    allTopicSubscriptions.forEach((subscription) => {
        subscription.unsubscribe();
        subscription.removeAllListeners();
    });
    allTopicSubscriptions.clear();

    // Clear messages when switching topics
    chatStore.update(state => ({ 
        ...state, 
        activeTopic: topic,
        messages: []
    }));

    // Add a small delay to allow Centrifuge to clean up internal state
    setTimeout(() => {
        if (topic === 'all') {
            // Subscribe to all individual topics for "all" view
            const individualTopics = ['general', 'tech', 'random'];
            
            individualTopics.forEach((individualTopic) => {
                const channelName = `chat:${individualTopic}`;
                
                try {
                    // Use getSubscription first to check if it exists, if not create new one
                    let subscription = centrifuge.getSubscription(channelName);
                    if (!subscription) {
                        subscription = centrifuge.newSubscription(channelName);
                    }

                    // Remove existing listeners to prevent duplicates
                    subscription.removeAllListeners();

                    subscription.on('publication', (ctx: any) => {
                        const message: Message = {
                            id: `${Date.now()}-${Math.random()}`,
                            username: ctx.data.username,
                            content: ctx.data.content,
                            timestamp: ctx.data.timestamp || Date.now(),
                            topic: individualTopic // Keep the original topic for display
                        };

                        chatStore.update(state => ({
                            ...state,
                            messages: [...state.messages, message]
                        }));
                    });

                    subscription.on('subscribed', () => {
                        console.log(`Subscribed to ${channelName} for "all" view`);
                    });

                    if (!subscription.subscribed) {
                        subscription.subscribe();
                    }
                    allTopicSubscriptions.set(individualTopic, subscription);
                } catch (error) {
                    console.error(`Error subscribing to ${channelName}:`, error);
                }
            });

            console.log('Subscribed to all channels for "all" topic view');
        } else {
            // Subscribe to specific topic
            const channelName = `chat:${topic}`;
            
            try {
                // Use getSubscription first to check if it exists, if not create new one
                let subscription = centrifuge.getSubscription(channelName);
                if (!subscription) {
                    subscription = centrifuge.newSubscription(channelName);
                }

                // Remove existing listeners to prevent duplicates
                subscription.removeAllListeners();

                subscription.on('publication', (ctx: any) => {
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

                subscription.on('subscribed', () => {
                    console.log(`Subscribed to ${channelName}`);
                });

                if (!subscription.subscribed) {
                    subscription.subscribe();
                }
                currentSubscription = subscription;
            } catch (error) {
                console.error(`Error subscribing to ${channelName}:`, error);
            }
        }
    }, 100); // Small delay to allow cleanup
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
    
    // Clean up all topic subscriptions
    allTopicSubscriptions.forEach((subscription) => {
        subscription.unsubscribe();
    });
    allTopicSubscriptions.clear();
    
    if (centrifuge) {
        centrifuge.disconnect();
        centrifuge = null;
    }
    
    chatStore.set(initialState);
};