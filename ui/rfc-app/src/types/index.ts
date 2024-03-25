export interface IRFC {
    id: string;
    title: string;
    authors: string;
    year: string;
    month: string;
    status: string;
    obsoleted_by: string;
    obsoletes: string;
    updated_by: string;
    updates: string;
    also: string;
  }
  
  export interface IRFCDetail {
    id: string;
    content: string;
  }
  