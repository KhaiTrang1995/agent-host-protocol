import XCTest
import AgentHostProtocol

final class ChatSourceTests: XCTestCase {

    func testChatSourceRoutesByKind() throws {
        let decoder = JSONDecoder()

        let fork = try decoder.decode(
            ChatSource.self,
            from: Data(#"{"kind":"fork","chat":"ahp-chat:/main","turnId":"turn-12"}"#.utf8)
        )
        let sideChat = try decoder.decode(
            ChatSource.self,
            from: Data(#"{"kind":"sideChat","chat":"ahp-chat:/main","turnId":"turn-active"}"#.utf8)
        )

        switch fork {
        case .fork(let value):
            XCTAssertEqual(value.kind, .fork)
            XCTAssertEqual(value.turnId, "turn-12")
        default:
            XCTFail("Expected fork variant")
        }

        switch sideChat {
        case .sideChat(let value):
            XCTAssertEqual(value.kind, .sideChat)
            XCTAssertEqual(value.turnId, "turn-active")
        default:
            XCTFail("Expected sideChat variant")
        }
    }

    func testChatSourceRejectsMissingOrUnknownKind() {
        let decoder = JSONDecoder()

        XCTAssertThrowsError(
            try decoder.decode(
                ChatSource.self,
                from: Data(#"{"chat":"ahp-chat:/main","turnId":"turn-12"}"#.utf8)
            )
        )
        XCTAssertThrowsError(
            try decoder.decode(
                ChatSource.self,
                from: Data(#"{"kind":"future","chat":"ahp-chat:/main","turnId":"turn-12"}"#.utf8)
            )
        )
    }
}
