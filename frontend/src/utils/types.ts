// Example types for your Go backend responses
export type AuthResponse = {
    token: string;
    user: {
      id: string;
      email: string;
    };
  };
  
  export type Competitor = {
    id: string;
    user_id: string;
    url: string;
    last_scraped_data: string;
    created_at: string;
  };

  export type User = {
    id: string;
    email: string;
    password: string;
    created_at: string;
  }

  export type PriceChange = {
    id: string;
    competitor_id: string;
    detected_change: string;
    ai_summary: string;
    created_at: string;
  }