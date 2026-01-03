import apiClient, { ApiResponse, handleApiError } from "./client";

/**
 * API client for chatbot-service
 */

export interface ChatMessage {
  role: "user" | "assistant";
  content: string;
}

export interface ChatContext {
  user_id?: string;
  current_step?: string;
  trip_search_id?: string;
  selected_trip?: {
    trip_id: string;
    origin: string;
    destination: string;
    departure_time: string;
    price: number;
  };
}

export interface ChatRequest {
  session_id?: string;
  message: string;
  history?: ChatMessage[];
  context?: ChatContext;
}

// Trip data structure returned from chatbot search
export interface ChatbotTripData {
  id: string;
  departure_time: string;
  arrival_time: string;
  origin: string;
  destination: string;
  price: number;
  available_seats: number;
  bus?: {
    name?: string;
    type?: string;
  };
  route?: {
    name?: string;
    estimated_duration?: number;
  };
}

export interface ChatResponse {
  message: string;
  intent?: string;
  action?: string;
  data?: {
    trips?: ChatbotTripData[];
    [key: string]: unknown;
  };
  context?: ChatContext;
  suggestions?: string[];
}

/**
 * Send a chat message to the chatbot
 */
export const sendChatMessage = async (
  req: ChatRequest,
): Promise<ChatResponse> => {
  try {
    const response = await apiClient.post<ApiResponse<ChatResponse>>(
      "/chatbot/api/v1/chat",
      req,
    );

    if (!response.data.data) {
      throw new Error("No response from chatbot");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Extract trip search parameters from natural language
 */
export const extractSearchParams = async (
  message: string,
): Promise<{
  origin: string;
  destination: string;
  departure_date?: string;
  passengers?: number;
}> => {
  try {
    const response = await apiClient.get<
      ApiResponse<{
        origin: string;
        destination: string;
        departure_date?: string;
        passengers?: number;
      }>
    >("/chatbot/api/v1/chat/extract-search", {
      params: { message },
    });

    if (!response.data.data) {
      throw new Error("Failed to extract search parameters");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};

/**
 * Get FAQ answer
 */
export const getFAQAnswer = async (question: string): Promise<string> => {
  try {
    const response = await apiClient.get<ApiResponse<string>>(
      "/chatbot/api/v1/chat/faq",
      {
        params: { question },
      },
    );

    if (!response.data.data) {
      throw new Error("No answer found");
    }

    return response.data.data;
  } catch (error) {
    const errorMessage = handleApiError(error);
    throw new Error(errorMessage);
  }
};
