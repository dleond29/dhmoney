import io.restassured.RestAssured;
import io.restassured.response.Response;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.junit.jupiter.api.DisplayName;
import org.springframework.core.annotation.Order;
import org.testng.annotations.Test;
import org.json.simple.JSONObject;

public class TestCasesLogout {

    private String baseUrl = "http://localhost:8080/api/users";

    @Test
    @Order(1)
    @DisplayName("logoutSuccessfully")
    public void logoutTestCase() throws ParseException {
        String login = "{\"email\":\"pepito1@gmail.com\"," +
                "\"password\":\"pepito\"}";

        RestAssured.baseURI = baseUrl;
        Response response = RestAssured
                .given().header("Content-Type", "application/json").body(login)
                .when().post("/login");

        JSONParser parser = new JSONParser();
        JSONObject jsonObject = (JSONObject) parser.parse(response.getBody().asString());
        String token = jsonObject.get("token").toString();

        RestAssured.baseURI = baseUrl;
        RestAssured
                .given().header("Content-Type", "application/json").header("Authorization", token)
                .when().post("/logout")
                .then().assertThat().statusCode(200);
    }


    @Test
    @Order(2)
    @DisplayName("forgotPasswordBadRequest")
    public void logoutWithoutTokenTestCase() {
        String token = "";

        RestAssured.baseURI = baseUrl;
        RestAssured
                .given().header("Content-Type", "application/json").header("Authorization", token)
                .when().post("/logout")
                .then().assertThat().statusCode(400);

    }
}
