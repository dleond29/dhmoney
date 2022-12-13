import com.fasterxml.jackson.databind.util.JSONPObject;
import io.restassured.RestAssured;
import org.junit.jupiter.api.DisplayName;
import org.springframework.core.annotation.Order;
import org.testng.annotations.Test;


public class TestCasesLogin {
    private String baseUrl = "http://localhost:8080/api/users";

    @Test
    @Order(1)
    @DisplayName("loginSuccessfully")
    public void LoginSuccessfullyTestCase() {
        String login = "{\"email\":\"pepito1@gmail.com\"," +
                "\"password\":\"pepito\"}";

        RestAssured.baseURI = baseUrl;
        RestAssured
                .given().header("Content-Type", "application/json").body(login)
                .when().post("/login")
                .then().assertThat().statusCode(200);
    }

    @Test
    @Order(2)
    @DisplayName("loginInvalidUserCredentials")
    public void LoginInvalidUserCredentialsTestCase() {
        String login = "{\"email\":\"noexiste@gmail.com\"," +
                "\"password\":\"pepito\"}";

        RestAssured.baseURI = baseUrl;
        RestAssured
                .given().header("Content-Type", "application/json").body(login)
                .when().post("/login")
                .then().assertThat().statusCode(400);
    }

    @Test
    @Order(3)
    @DisplayName("loginAllFieldsAreRequired")
    public void LoginAllFieldsAreRequiredTestCase() {
        String login = "{\"email\":\"noexiste@gmail.com\"}";

        RestAssured.baseURI = baseUrl;
        RestAssured
                .given().header("Content-Type", "application/json").body(login)
                .when().post("/login")
                .then().assertThat().statusCode(400);
    }

    @Test
    @Order(3)
    @DisplayName("loginEmailNotVerified")
    public void LoginEmailNotVerifiedTestCase() {
        String login = "{\"email\":\"noverified@gmail.com\"," +
                "\"password\":\"noverified\"}";

        RestAssured.baseURI = baseUrl;
        RestAssured
                .given().header("Content-Type", "application/json").body(login)
                .when().post("/login")
                .then().assertThat().statusCode(401);
    }

    @Test
    @Order(4)
    @DisplayName("loginUserNotExists")
    public void LoginUserNotExistsTestCase() {
        String login = "{\"email\":\"notfound@gmail.com\"," +
                "\"password\":\"notfoun\"}";

        RestAssured.baseURI = baseUrl;
        RestAssured
                .given().header("Content-Type", "application/json").body(login)
                .when().post("/login")
                .then().assertThat().statusCode(404);
    }

}
