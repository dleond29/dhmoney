import io.restassured.RestAssured;
import org.junit.jupiter.api.DisplayName;
import org.springframework.core.annotation.Order;
import org.testng.annotations.Test;

public class TestCasesForgotPassword {

    private String baseUrl = "http://localhost:8080/api/users";

    @Test
    @Order(1)
    @DisplayName("forgotPasswordSuccessfully")
    public void forgotPasswordSuccessfullyTestCase() {
        String forgotPassword = "{\"email\":\"pepito1@gmail.com\"}";

        RestAssured.baseURI = baseUrl;
        RestAssured
                .given().header("Content-Type", "application/json").body(forgotPassword)
                .when().post("/forgot")
                .then().assertThat().statusCode(200);
    }

    @Test
    @Order(2)
    @DisplayName("forgotPasswordBadRequest")
    public void forgotPasswordBadRequestTestCase() {
        String forgotPassword = "";

        RestAssured.baseURI = baseUrl;
        RestAssured
                .given().header("Content-Type", "application/json").body(forgotPassword)
                .when().post("/forgot")
                .then().assertThat().statusCode(400);
    }
}
