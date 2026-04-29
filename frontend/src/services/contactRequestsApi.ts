import { api } from './http';
import type { ApiResponse, BackendReceivedContactRequest } from './types';

export interface ReceivedContactRequest {
  id: number;
  listingId: number;
  listingTitle: string;
  buyerId: number;
  buyerName: string;
  buyerEmail: string;
  message: string;
  status: string;
  createdAt: string;
}

function mapReceivedContactRequest(
  request: BackendReceivedContactRequest
): ReceivedContactRequest {
  return {
    id: request.id,
    listingId: request.listing_id,
    listingTitle: request.listing_title,
    buyerId: request.buyer_id,
    buyerName: request.buyer_name,
    buyerEmail: request.buyer_email,
    message: request.message,
    status: request.status,
    createdAt: request.created_at,
  };
}

export const contactRequestsApi = {
  getReceived: async (): Promise<{ data: ApiResponse<ReceivedContactRequest[]> }> => {
    const response = await api.get<ApiResponse<BackendReceivedContactRequest[]>>(
      '/contact-requests/received'
    );

    if (response.data.success && response.data.data) {
      response.data.data = response.data.data.map(mapReceivedContactRequest) as never;
    }

    return response as unknown as { data: ApiResponse<ReceivedContactRequest[]> };
  },
};
