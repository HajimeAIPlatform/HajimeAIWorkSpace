from pythonp.apps.data.utils.utils import remove_non_utf8

def test_remove_non_utf8():
    text = "Hello, 你好, こんにちは, 안녕하세요"
    text2 = "第八十八签解签：\n\n张子房误中副车  中平\n"
    cleaned_text = remove_non_utf8(text)
    cleaned_text2 = remove_non_utf8(text2)
    assert cleaned_text2 == "第八十八签解签：\n\n张子房误中副车  中平\n"
    assert cleaned_text == "Hello, 你好, こんにちは, 안녕하세요"


def test_smoke_test():
    assert True

if __name__ == "__main__":
    import pytest
    raise SystemExit(pytest.main([__file__]))