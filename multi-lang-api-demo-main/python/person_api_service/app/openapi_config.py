from fastapi.openapi.utils import get_openapi

def custom_openapi(app):
    if app.openapi_schema:
        return app.openapi_schema
    openapi_schema = get_openapi(
        title="Person API Service",
        version="1.0.0",
        description="FastAPI REST API for Person management with SQLite and MongoDB",
        routes=app.routes,
    )
    app.openapi_schema = openapi_schema
    return app.openapi_schema
