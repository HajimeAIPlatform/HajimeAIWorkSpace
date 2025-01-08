from typing import Dict, List, Optional, Any

from pydantic import BaseModel, Field

from pythonp.apps.tokenfate.dify_client.models.base import ResponseMode, File
from pythonp.apps.tokenfate.dify_client.models.completion import CompletionResponse


class ChatRequest(BaseModel):
    query: str
    inputs: Dict[str, Any] = Field(default_factory=dict)
    response_mode: ResponseMode
    user: str
    conversation_id: Optional[str] = ""
    files: List[File] = []
    auto_generate_name: bool = True


class ChatResponse(CompletionResponse):
    pass


class ChatSuggestRequest(BaseModel):
    user: str


class ChatSuggestResponse(BaseModel):
    result: str
    data: List[str] = []
