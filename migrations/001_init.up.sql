CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY,
    
    creator_id UUID NOT NULL,
    creator_phone_number varchar(64) NOT NULL,
    creator_name varchar(256) NOT NULL,
    message TEXT,
   
    assigned_to_id UUID,
    is_resolved BOOLEAN DEFAULT false,
    resolved_at TIMESTAMP,
    
    with_success BOOLEAN,
    taken_action TEXT,
    customer_rating INTEGER CHECK (customer_rating>=1 AND customer_rating<=5),
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)